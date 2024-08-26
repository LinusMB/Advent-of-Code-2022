package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
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

	reg := 1
	regs := []int{}
	cs := []int{20, 60, 100, 140, 180, 220}

	sc := bufio.NewScanner(file)
	for sc.Scan() {
		line := sc.Text()
		switch {
		case strings.HasPrefix(line, "addx"):
			var n int
			fmt.Sscanf(line, "addx %d", &n)
			regs = append(regs, reg, reg)
			reg = reg + n
		case strings.HasPrefix(line, "noop"):
			regs = append(regs, reg)
		default:
			panic("Unexpected line")
		}
	}
	sum := 0
	for _, c := range cs {
		sum += regs[c-1] * c
	}
	fmt.Printf("Solution Part 1: %v\n", sum)

	fmt.Printf("Solution Part 2\n")
	for i := range regs {
		p := i % 40

		if (p >= regs[i]-1) && (p <= regs[i]+1) {
			print("#")
		} else {
			print(" ")
		}
		if (i+1)%40 == 0 {
			print("\n")
		}

	}
}
