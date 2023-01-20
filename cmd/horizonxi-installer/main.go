package main

import (
	"bytes"
	"crypto/sha256"
	_ "embed"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"time"

	"github.com/anacrolix/torrent"
	grab "github.com/cavaliergopher/grab/v3"
)

const (
	LOGO = `
	 _   _            _               __  _____ 
	| | | | ___  _ __(_)_______  _ __ \ \/ /_ _|
	| |_| |/ _ \| '__| |_  / _ \| '_ \ \  / | | 
	|  _  | (_) | |  | |/ / (_) | | | |/  \ | | 
	|_| |_|\___/|_|  |_/___\___/|_| |_/_/\_\___|
	`

	HorizonXIZipMagnet = "magnet:?xt=urn:btih:4eecae8431428820347314bc002492e210f29612&dn=HorizonXI.zip&tr=udp%3A%2F%2Fopentracker.i2p.rocks%3A6969%2Fannounce&tr=https%3A%2F%2Ftracker.nanoha.org%3A443%2Fannounce&tr=https%3A%2F%2Ftracker.lilithraws.org%3A443%2Fannounce&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337%2Fannounce&tr=udp%3A%2F%2Ftracker.openbittorrent.com%3A6969%2Fannounce&tr=https%3A%2F%2Fopentracker.i2p.rocks%3A443%2Fannounce&tr=udp%3A%2F%2Ftracker1.bt.moack.co.kr%3A80%2Fannounce&tr=udp%3A%2F%2Ftracker.torrent.eu.org%3A451%2Fannounce&tr=udp%3A%2F%2Ftracker.tiny-vps.com%3A6969%2Fannounce&tr=udp%3A%2F%2Fpublic.tracker.vraphim.com%3A6969%2Fannounce&tr=udp%3A%2F%2Fp4p.arenabg.com%3A1337%2Fannounce&tr=udp%3A%2F%2Fopen.stealth.si%3A80%2Fannounce&tr=udp%3A%2F%2Fopen.demonii.com%3A1337%2Fannounce&tr=udp%3A%2F%2Fmovies.zsw.ca%3A6969%2Fannounce&tr=udp%3A%2F%2Fipv4.tracker.harry.lu%3A80%2Fannounce&tr=udp%3A%2F%2Fexplodie.org%3A6969%2Fannounce&tr=udp%3A%2F%2Fexodus.desync.com%3A6969%2Fannounce&tr=udp%3A%2F%2F9.rarbg.com%3A2810%2Fannounce&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337&tr=udp%3A%2F%2Fexplodie.org%3A6969"
	SHA256Sum          = "af374c21eda3dc7517fc9e5b0da46eba9a928574a3c6792ea74de590571874c2"
	DataFile           = "HorizonXI.zip"
	WineURL            = "https://github.com/GloriousEggroll/wine-ge-custom/releases/download/GE-Proton7-35/wine-lutris-GE-Proton7-35-x86_64.tar.xz"
	DgVoodoo2URL       = "http://dege.freeweb.hu/dgVoodoo2/bin/dgVoodoo2_79_3.zip"
	DXVKURL            = "https://github.com/doitsujin/dxvk/releases/download/v2.0/dxvk-2.0.tar.gz"

	ResetCursor = "\033[1F"

	magic    = 0x5a4d
	e_lfanew = 0x3c
)

//go:embed horizonxi
var LauncherScript []byte

//go:embed dgVoodoo.conf
var DgVoodoo2Config []byte

//go:embed hkey
var HkeyContents []byte

var WinePrefix = ""

var (
	installPath      = flag.String("p", os.Getenv("HOME")+"/HorizonXI", "install path")
	winePath         = path.Join(*installPath, "lutris-GE-Proton7-35-x86_64/bin/wine")
	dataFile         = flag.String("d", "", "path to HorizonXI.zip (optional, to skip download)")
	RequiredPrograms = [...]string{"tar", "xz", "unzip"}
)

func main() {
	flag.Parse()
	fmt.Println(LOGO)

	// check if dependencies are installed
	for _, program := range RequiredPrograms {
		if !CheckForProgram(program) {
			log.Fatalf("\"%s\" is a required program -- install it and then run the installer again", program)
		}
	}

	// check if install dir exists
	if _, err := os.Stat(*installPath); err != nil {
		os.Mkdir(*installPath, 0755)
	}

	WinePrefix = path.Join(*installPath, ".wine")

	// check if using existing data file
	if *dataFile != "" {
		log.Print("copying data file into install directory")
		if err := exec.Command("cp", *dataFile, *installPath).Run(); err != nil {
			log.Fatalf("could not copy data file: %s", err)
		}
	}

	// change to install directory
	if err := os.Chdir(*installPath); err != nil {
		log.Fatal("cannot access install path, please check permissions on: " + *installPath)
	}

	// obtain data files - if they already exist, check integrity
	// and download the data files again if integrity check fails
	if _, err := os.Stat(DataFile); err != nil {
		DownloadDataFiles()
	} else {
		log.Print("verifying files")
		if !CheckFileHash(DataFile, SHA256Sum) {
			log.Print("verify failed, downloading...")
			DownloadDataFiles()
		}
	}
	InstallDataFiles()
	InstallWine()

	InstallDgVoodoo2()
	InstallDXVK()
	LargeAddressPatch(path.Join(*installPath, "HorizonXI/bootloader/horizon-loader.exe"))
	Enable60FPS()
	InstallLauncher()
	log.Printf("install complete -- to play, run %s/horizonxi", *installPath)
}

func Enable60FPS() {
	log.Print("enabling 60fps")
	f, err := os.OpenFile(path.Join(*installPath, "HorizonXI/scripts/default.txt"), os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	f.WriteString("\n/fps 1\n")
}

func InstallDXVK() {
	log.Print("downloading DXVK")
	resp, err := grab.Get(".", DXVKURL)
	if err != nil {
		log.Fatalf("could not download DXVK: %s", err)
	}
	log.Print("extracting DXVK")
	cmd := exec.Command("tar", "xzf", resp.Filename)
	if err := cmd.Run(); err != nil {
		log.Fatalf("error extracting DXVK: %s", err)
	}
	log.Print("installing DXVK")
	cmd = exec.Command("./setup_dxvk.sh", "install")
	cmd.Dir = path.Join(*installPath, "dxvk-2.0")
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("WINEPREFIX=%s", WinePrefix))

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Print(string(output))
		log.Fatalf("error installing dxvk (output dumped): %s", err)
	}
}

func InstallDgVoodoo2() {
	log.Print("downloading dgVoodoo2")
	resp, err := grab.Get(".", DgVoodoo2URL)
	if err != nil {
		log.Fatalf("could not download dgVoodoo2: %s", err)
	}
	log.Print("extracting dgVoodoo2")
	cmd := exec.Command("unzip", "-u", resp.Filename, "-d", "dgVoodoo2")
	if err := cmd.Run(); err != nil {
		log.Fatalf("error extracting dgVoodoo2: %s", err)
	}
	os.Chdir(*installPath)
	for _, dll := range []string{"D3D8.dll", "D3D9.dll", "D3DImm.dll", "DDraw.dll"} {
		err = os.Rename(path.Join("dgVoodoo2/MS/x86", dll), path.Join("HorizonXI/bootloader", dll))
		if err != nil {
			log.Fatalf("error install dgVoodoo2 file %s: %s", dll, err)
		}
	}

	err = os.WriteFile(path.Join(*installPath, "hkey"), HkeyContents, 0755)
	if err != nil {
		log.Fatalf("error writiying hkey tool: %s", err)
	}

	cmd = exec.Command("./hkey")
	cmd.Dir = *installPath
	if err := cmd.Run(); err != nil {
		log.Fatalf("error setting DllOverrides: %s", err)
	}
	err = os.WriteFile("HorizonXI/bootloader/dgVoodoo.conf", DgVoodoo2Config, 0644)
	if err != nil {
		log.Fatalf("error writing dgVoodoo2 config: %s", err)
	}
}

func InstallLauncher() {
	os.WriteFile("horizonxi", LauncherScript, 0755)
}

func InstallDataFiles() {
	log.Print("installing game files")
	cmd := exec.Command("unzip", DataFile)
	if err := cmd.Run(); err != nil {
		log.Fatalf("error installing game files: %s", err)
	}
}

func InstallWine() {
	log.Print("downloading wine")
	resp, err := grab.Get(".", WineURL)
	if err != nil {
		log.Fatalf("could not download wine: %s", err)
	}
	log.Print("installing wine")
	cmd := exec.Command("tar", "xJf", resp.Filename)
	if err := cmd.Run(); err != nil {
		log.Fatalf("error installing wine: %s", err)
	}
	winecfg := path.Join(*installPath, "lutris-GE-Proton7-35-x86_64", "bin", "wineboot")
	cmd = exec.Command(winecfg)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("WINEPREFIX=%s", WinePrefix))
	cmd.Env = append(cmd.Env, "WINEARCH=win32")
	if err := cmd.Run(); err != nil {
		log.Fatalf("error running winecfg: %s", err)
	}
}

func CheckFileHash(file, hash string) bool {
	f, err := os.Open(file)
	if err != nil {
		log.Fatalf("could not open %s: %s", file, err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatalf("input/output error while checking file integrity: %s", err)
	}

	return fmt.Sprintf("%x", h.Sum(nil)) == hash
}

func DownloadDataFiles() {
	log.Print("initiating download")
	c, _ := torrent.NewClient(nil)
	defer c.Close()
	t, _ := c.AddMagnet(HorizonXIZipMagnet)
	<-t.GotInfo()
	t.DownloadAll()
	for {
		stats := t.Stats()
		percentComplete := float64(stats.PiecesComplete) / float64(t.NumPieces()) * 100.0
		log.Printf("downloading game files: %.2f%%\n", percentComplete)

		if percentComplete >= 100.0 {
			break
		}

		time.Sleep(1 * time.Second)

		// reset cursor to beginning of previous line
		fmt.Printf(ResetCursor)
	}
	c.WaitAll()
	log.Print("data download complete")
}

func CheckForProgram(program string) bool {
	cmd := exec.Command("which", program)
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

func LargeAddressPatch(inputFile string) {
	log.Print("patching HorizonXI bootloader to support large addresses")
	flag.Parse()
	f, err := os.OpenFile(inputFile, os.O_RDWR, 0755)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	buf := make([]byte, 4)

	// read MZ header
	io.ReadAtLeast(f, buf, 2)
	header := binary.LittleEndian.Uint16(buf[0:2])
	if header != magic {
		log.Fatalf("not a valid MS-DOS executable")
	}

	// seek to where the PE address is stored
	f.Seek(e_lfanew, io.SeekStart)

	// read the PE address
	io.ReadAtLeast(f, buf, 4)
	peAddress := int64(binary.LittleEndian.Uint32(buf[0:4]))

	// seek to PE
	f.Seek(peAddress, 0)

	io.ReadAtLeast(f, buf, 2)
	if !bytes.Equal(buf[0:2], []byte{0x50, 0x45}) {
		log.Fatalf("invalid PE file")
	}
	f.Seek(peAddress+0x16, 0)
	io.ReadAtLeast(f, buf, 2)
	flags := binary.LittleEndian.Uint16(buf[0:2])
	flags |= 0x20
	f.Seek(peAddress+0x16, 0)
	binary.Write(f, binary.LittleEndian, flags)
	log.Println("patch succeeded!")
}
