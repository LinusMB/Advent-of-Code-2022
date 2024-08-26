package main

import (
	"cmp"
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

func max[T cmp.Ordered](x, y T) T {
	if x > y {
		return x
	}
	return y
}

func main() {
	file, err := os.Open("input.txt")
	die(err, "Open file")
	defer file.Close()

	bytes, err := io.ReadAll(file)
	die(err, "Read file")

	var grid []uint8
	var ncols, nrows int

	for i := range bytes {
		if bytes[i] == '\n' {
			nrows++
			if ncols == 0 {
				ncols = i
			}
			continue
		}
		grid = append(grid, uint8(bytes[i])-uint8('0'))
	}

	edge := func(i int) bool {
		r := i / ncols
		c := i % ncols
		return r == 0 || r == nrows-1 || c == 0 || c == ncols-1
	}

	visible := func(i int) bool {
		for _, j := range []int{-1, +1, -ncols, +ncols} {
			for k := i; k == i || grid[i] > grid[k]; k += j {
				if edge(k) {
					return true
				}
			}
		}
		return false
	}

	score := func(i int) int {
		s := 1
		for _, j := range []int{-1, +1, -ncols, +ncols} {
			d := 0
			for k := i; !edge(k) && (k == i || grid[i] > grid[k]); k += j {
				d++
			}
			s *= d
		}
		return s
	}

	sol1, sol2 := 0, 0
	for i := 0; i < nrows*ncols; i++ {
		if visible(i) {
			sol1++
		}
		sol2 = max(sol2, score(i))
	}

	fmt.Printf("Solution Part 1: %d\n", sol1)
	fmt.Printf("Solution Part 2: %d\n", sol2)
}
