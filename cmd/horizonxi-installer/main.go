package main

import (
	"crypto/sha256"
	_ "embed"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
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

	WineURL = "https://github.com/GloriousEggroll/wine-ge-custom/releases/download/GE-Proton7-35/wine-lutris-GE-Proton7-35-x86_64.tar.xz"

	ResetCursor = "\033[1F"
)

//go:embed horizonxi
var LauncherScript []byte

var (
	installPath      = flag.String("p", os.Getenv("HOME")+"/HorizonXI", "install path")
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

	// change to install directory
	if err := os.Chdir(*installPath); err != nil {
		log.Fatal("cannot access install path, please check permissions on: " + *installPath)
	}

	// obtain data files - if they already exist, check integrity
	// and download the data files again if integrity check fails
	if _, err := os.Stat(DataFile); err != nil {
		DownloadDataFiles()
	} else {
		log.Print("data file detected, checking integrity")
		if !CheckFileHash(DataFile, SHA256Sum) {
			log.Print("data file integrity check failed")
			DownloadDataFiles()
		}
	}
	InstallDataFiles()
	InstallWine()
	InstallLauncher()
	log.Printf("install complete -- to play, run %s/horizonxi", *installPath)
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
