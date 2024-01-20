package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
)

// Helper functions

func countBytesInFile(filename string) ([]uint64, uint64, error) {
	// Open file
	file, err := os.Open(filename)
	if err != nil {
		return nil, 0, err
	}
	defer file.Close()

	byteValueCounter := make([]uint64, 256)
	byteCount := uint64(0)

	// Read file in chunks
	buffer := make([]byte, 1024*256)
	for {
		bytesRead, err := file.Read(buffer)
		// Stop when end of file reached
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, 0, err
		}
		// Count occurances of each byte value
		if bytesRead > 0 {
			for i := 0; i < bytesRead; i++ {
				byteValueCounter[int(buffer[i])]++
			}
			byteCount += uint64(bytesRead)
		}
	}

	return byteValueCounter, byteCount, nil
}

func calculateEntropyBits(values []uint64, count uint64) float64 {
	// Calculate bits (shannons) of entropy
	// https://en.wikipedia.org/wiki/Entropy_(information_theory)
	var entropy float64 = 0.0
	for _, v := range values {
		if v > 0 {
			px := float64(v) / float64(count)
			entropy += px * math.Log2(px)
		}
	}
	return -entropy
}

func checkFileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !(errors.Is(err, os.ErrNotExist) || errors.Is(err, os.ErrPermission))
}

// Main program

func main() {
	// Prepare for parsing command line args
	showGraph := flag.Bool("g", false, "Graph output")
	numberOnlyOutput := flag.Bool("q", false, "Number-only output (quiet mode)")

	_, _ = numberOnlyOutput, showGraph

	// Custom Usage message
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "%s [-q/-g] file\n", os.Args[0])
		flag.PrintDefaults()
	}

	// Parse args
	flag.Parse()

	// Check if user provided at least 1 argument
	if len(os.Args) < 2 {
		fmt.Fprintf(flag.CommandLine.Output(), "file not provided\n")
		flag.Usage()
		os.Exit(2)
	}

	// Get filename (last arg)
	filename := os.Args[len(os.Args)-1]

	// Check if file exists
	if !checkFileExists(filename) {
		fmt.Fprintf(flag.CommandLine.Output(), "file %s does not exist or has wrong permissions set\n", filename)
		flag.Usage()
		os.Exit(1)
	}

	// Count byte values in file
	byteValueCounter, byteCount, err := countBytesInFile(filename)
	if err != nil {
		fmt.Fprintf(flag.CommandLine.Output(), "%s\n", err.Error())
		flag.Usage()
		os.Exit(1)
	}

	// Calculate entropy in bits (might add different types later)
	entropy := calculateEntropyBits(byteValueCounter, byteCount)

	// Display output
	// TODO: different display types
	fmt.Printf("Entropy: %.06f\n", entropy)
}
