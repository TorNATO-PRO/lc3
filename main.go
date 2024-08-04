package main

import (
	"bytes"
	"encoding/binary"
	"lc3/internal/constants"
	"lc3/pkg/cpu"
	"log"
	"math"
	"os"
)

func readImage(filename string) ([constants.MemoryMax]uint16, error) {
	m := [constants.MemoryMax]uint16{}

	file, err := os.Open(filename)

	if err != nil {
		return m, err
	}

	defer file.Close()

	stats, err := file.Stat()
	if err != nil {
		return m, err
	}

	// load the origin
	var origin uint16

	headerBytes := make([]byte, 2)
	_, err = file.Read(headerBytes)
	if err != nil {
		return m, err
	}

	headerBuffer := bytes.NewBuffer(headerBytes)
	err = binary.Read(headerBuffer, binary.BigEndian, &origin)
	if err != nil {
		return m, err
	}

	log.Printf("Origin memory location: 0x%04X", origin)
	size := stats.Size()
	byteArr := make([]byte, size)

	log.Printf("Creating memory buffer: %d bytes", size)

	_, err = file.Read(byteArr)
	if err != nil {
		return m, err
	}

	buffer := bytes.NewBuffer(byteArr)

	for i := origin; i < math.MaxUint16; i++ {
		var val uint16
		binary.Read(buffer, binary.BigEndian, &val)
		m[i] = val
	}

	return m, err
}

func loadArguments() [][constants.MemoryMax]uint16 {
	args := os.Args[1:]

	if len(args) < 1 {
		log.Fatal("lc3 [image-file1] ...\n")
	}

	var images [][65536]uint16

	for _, arg := range args {
		image, err := readImage(arg)

		if err != nil {
			log.Fatalf("failed to load image: %s, %v", arg, err)
		}

		images = append(images, image)
	}

	return images
}

func main() {
	args := loadArguments()

	for _, args := range args {
		cpu := cpu.NewCPU()

		err := cpu.Run(args)

		if err != nil {
			log.Fatalf("Execution failed %v", err)
		}
	}
}
