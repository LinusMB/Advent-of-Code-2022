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

func priority(c rune) int {
	if c >= 'a' && c <= 'z' {
		return int(c) - int('a') + 1
	}
	if c >= 'A' && c <= 'Z' {
		return int(c) - int('A') + 27
	}
	return 0
}

func main() {
	file, err := os.Open("input.txt")
	die(err, "Open file")
	defer file.Close()

	var lines []string
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}

	sum1, sum2 := 0, 0

	makeCharSet := func(s string) map[rune]bool {
		set := make(map[rune]bool)
		for _, c := range s {
			set[c] = true
		}
		return set
	}

	for i := range lines {
		n := len(lines[i]) / 2
		left, right := lines[i][:n], lines[i][n:]
		charSets := make([]map[rune]bool, 2)
		charSets[0] = makeCharSet(left)
		charSets[1] = makeCharSet(right)

		var intsct rune
		for c := range charSets[0] {
			if _, inSecond := charSets[1][c]; inSecond {
				intsct = c
				break
			}
		}
		sum1 += priority(intsct)
	}

	for i := 0; i < len(lines); i += 3 {
		charSets := make([]map[rune]bool, 3)
		for j := 0; j < 3; j++ {
			charSets[j] = makeCharSet(lines[i+j])
		}
		var intsct rune
		for c := range charSets[0] {
			_, inSecond := charSets[1][c]
			_, inThird := charSets[2][c]
			if inSecond && inThird {
				intsct = c
				break
			}
		}
		sum2 += priority(intsct)
	}

	fmt.Printf("Solution Part 1: %d\n", sum1)
	fmt.Printf("Solution Part 2: %d\n", sum2)
}
