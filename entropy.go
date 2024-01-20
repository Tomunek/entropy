package main

import (
	"fmt"
	"io"
	"math"
	"os"
)

func countBytesInFile(file os.File) ([]uint64, uint64) {
	byteValueCounter := make([]uint64, 256)
	byteCount := uint64(0)

	buffer := make([]byte, 1024*256)
	for {
		bytesRead, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s", err.Error())
			os.Exit(2)
		}
		if bytesRead > 0 {
			for i := 0; i < bytesRead; i++ {
				byteValueCounter[int(buffer[i])]++
			}
			byteCount += uint64(bytesRead)
		}
	}

	return byteValueCounter, byteCount
}

func calculateEntropy(values []uint64, count uint64) float64 {
	var entropy float64 = 0.0
	for _, v := range values {
		if v > 0 {
			px := float64(v) / float64(count)
			entropy += px * math.Log2(px)
		}
	}
	return -entropy
}

func main() {
	var filename string = ""
	file, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s", err.Error())
		os.Exit(1)
	}
	defer file.Close()

	byteValueCounter, byteCount := countBytesInFile(*file)

	entropy := calculateEntropy(byteValueCounter, byteCount)

	fmt.Printf("Entropy: %.06f\n", entropy)
}
