package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
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
		// Count occurrences of each byte value
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
	graphBins := flag.Int("b", 10, "Number of bins on graph")
	graphLength := flag.Int("l", 30, "Max graph bar length")
	showGraph := flag.Bool("g", false, "Show output graph ")
	numberOnlyOutput := flag.Bool("q", false, "Number-only output (quiet mode)")

	// Custom Usage message
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "%s [flags] file\n", os.Args[0])
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
	// Trick to prevent ugly -0.0 from displaying
	// At this point, entropy can be assumed to be >= 0
	if math.Signbit(entropy) {
		entropy = -entropy
	}

	// Display output
	if *numberOnlyOutput {
		// Number only (quiet) mode
		fmt.Println(entropy)
	} else {
		// Normal mode
		fmt.Printf("Filename: %s\n", filename)
		fmt.Printf("Bytes: %d\n", byteCount)
		fmt.Printf("Entropy: %.06f Sh\n", math.Round(entropy*1000000.0)/1000000.0)
	}

	if *showGraph {
		graphChar := "█"
		binsCount := *graphBins
		maxDisplayLength := *graphLength
		bins := make([]uint64, binsCount)
		maxBinValue := uint64(0)

		// Divide values into bins
		for i := 0; i < binsCount; i++ {
			binStart := 256 * i / binsCount
			binEnd := 256 * (i + 1) / binsCount
			for j := binStart; j < binEnd; j++ {
				bins[i] += byteValueCounter[j]
			}
			if maxBinValue < bins[i] {
				maxBinValue = bins[i]
			}
		}

		// Display graph
		displayLengthMultiplier := float64(maxDisplayLength) / float64(maxBinValue)
		fmt.Printf("\nValues     0 Counts%s%d\n", strings.Repeat(" ", maxDisplayLength-11), maxBinValue)
		for i := 0; i < binsCount; i++ {
			fmt.Printf("%3d - %3d: ", 256*i/binsCount, (256*(i+1)/binsCount)-1)
			fmt.Printf("%s\n", strings.Repeat(graphChar, int(math.Round(float64(bins[i])*displayLengthMultiplier))))
		}
	}
}
