package main

import (
	"bufio"
	"fmt"
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

	sc := bufio.NewScanner(file)

	count1, count2 := 0, 0

	for sc.Scan() {
		line := sc.Text()
		if line == "" {
			continue
		}
		var iv1, iv2 struct{ l, r int }
		fmt.Sscanf(line, "%d-%d,%d-%d", &iv1.l, &iv1.r, &iv2.l, &iv2.r)

		if (iv1.l <= iv2.l && iv1.r >= iv2.r) || (iv2.l <= iv1.l && iv2.r >= iv1.r) {
			count1++
		}

		if (iv1.r >= iv2.l && iv1.l <= iv2.r) || (iv1.l <= iv2.l && iv1.r >= iv2.r) {
			count2++
		}
	}
	fmt.Printf("Solution Part 1: %d\n", count1)
	fmt.Printf("Solution Part 2: %d\n", count2)
}
