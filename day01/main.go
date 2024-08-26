package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
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

	caloriesAcc := 0
	caloriesMax := [3]int{}

	updateMax := func(calories int) {
		for i := 0; i < 3; i++ {
			if calories <= caloriesMax[i] {
				continue
			}
			for j := 2; j > i; j-- {
				caloriesMax[j] = caloriesMax[j-1]
			}
			caloriesMax[i] = calories
			break
		}
	}

	for sc.Scan() {
		if line := sc.Text(); line != "" {
			calories, err := strconv.Atoi(line)
			die(err, "Parse line as number")
			caloriesAcc += calories
			continue
		}
		updateMax(caloriesAcc)
		caloriesAcc = 0
	}
	updateMax(caloriesAcc)

	fmt.Printf("Solution Part 1: %d\n", caloriesMax[0])
	fmt.Printf("Solution Part 2: %d\n", caloriesMax[0]+caloriesMax[1]+caloriesMax[2])
}
