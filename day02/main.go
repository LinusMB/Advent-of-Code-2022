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

	scoreTable := map[string]struct{ f, s int }{
		"A X": {3 + 1, 0 + 3},
		"B X": {0 + 1, 0 + 1},
		"C X": {6 + 1, 0 + 2},
		"A Y": {6 + 2, 3 + 1},
		"B Y": {3 + 2, 3 + 2},
		"C Y": {0 + 2, 3 + 3},
		"A Z": {0 + 3, 6 + 2},
		"B Z": {6 + 3, 6 + 3},
		"C Z": {3 + 3, 6 + 1},
	}

	score1, score2 := 0, 0

	sc := bufio.NewScanner(file)

	for sc.Scan() {
		line := sc.Text()
		score1 += scoreTable[line].f
		score2 += scoreTable[line].s
	}

	fmt.Printf("Solution Part 1: %d\n", score1)
	fmt.Printf("Solution Part 2: %d\n", score2)
}
