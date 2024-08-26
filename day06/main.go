package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

func die(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %v", msg, err)
	}
}

func main() {
	file, err := os.Open("input.txt")
	die(err, "Open file")
	defer file.Close()

	bytes, err := io.ReadAll(file)
	die(err, "Read file")

	unique := func(bs []byte) bool {
		seen := make(map[byte]bool)
		for _, b := range bs {
			if seen[b] {
				return false
			}
			seen[b] = true
		}
		return true
	}

	marker := func(nchars int) int {
		for i := 0; i <= len(bytes)-nchars; i++ {
			if unique(bytes[i : i+nchars]) {
				return i + nchars
			}
		}
		return -1
	}

	fmt.Printf("Solution Part 1: %d\n", marker(4))
	fmt.Printf("Solution Part 2: %d\n", marker(14))
}
