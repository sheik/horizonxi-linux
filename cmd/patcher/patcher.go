package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

// const MZHeader = []byte{0x4d, 0x5a}
const magic = 0x5a4d
const e_lfanew = 0x3c

var (
	inputFile = flag.String("i", "", "input file")
)

func main() {
	flag.Parse()
	f, err := os.OpenFile(*inputFile, os.O_RDWR, 0755)
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
	fmt.Printf("Flags: %x\n", flags)
	flags |= 0x20
	f.Seek(peAddress+0x16, 0)
	fmt.Printf("Flags After: %x\n", flags)
	binary.Write(f, binary.LittleEndian, flags)
	fmt.Println("patch succeeded!")
}
