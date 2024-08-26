package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"golang.org/x/exp/constraints"
)

func die(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %v", msg, err)
	}
}

func assert(b bool, msg string) {
	if !b {
		panic(msg)
	}
}

func sign[T constraints.Signed | constraints.Float](x T) T {
	if x < 0 {
		return -1
	}
	return 1
}

func abs[T constraints.Signed | constraints.Float](x T) T {
	return sign(x) * x
}

type Vec struct{ x, y int }

type Rope struct {
	knots []Vec
	seen  map[Vec]bool
}

func newRope(s Vec, n int) *Rope {
	r := Rope{
		knots: make([]Vec, n),
		seen:  make(map[Vec]bool),
	}
	for i := range r.knots {
		r.knots[i] = s
	}
	r.seen[s] = true

	return &r
}

func (r *Rope) move(dv Vec) {
	vecAdd := func(a, b Vec) Vec { return Vec{a.x + b.x, a.y + b.y} }
	vecSub := func(a, b Vec) Vec { return Vec{a.x - b.x, a.y - b.y} }

	r.knots[0] = vecAdd(r.knots[0], dv)

	for i := 1; i < len(r.knots); i++ {
		dv = vecSub(r.knots[i-1], r.knots[i])
		if abs(dv.x) <= 1 && abs(dv.y) <= 1 {
			return
		}
		assert(abs(dv.x) <= 2 && abs(dv.y) <= 2, "Unexpected difference vector")
		if abs(dv.x) > 1 {
			dv.x = sign(dv.x)
		}
		if abs(dv.y) > 1 {
			dv.y = sign(dv.y)
		}
		r.knots[i] = vecAdd(r.knots[i], dv)
	}

	r.seen[r.knots[len(r.knots)-1]] = true
}

func main() {
	file, err := os.Open("input.txt")
	die(err, "Open file")
	defer file.Close()

	rope1 := newRope(Vec{0, 0}, 2)
	rope2 := newRope(Vec{0, 0}, 10)

	sc := bufio.NewScanner(file)
	for sc.Scan() {
		line := sc.Text()
		var (
			inst byte
			n    int
		)
		fmt.Sscanf(line, "%c %d", &inst, &n)

		dv, ok := map[byte]Vec{
			'R': Vec{1, 0}, 'L': Vec{-1, 0}, 'U': Vec{0, 1}, 'D': Vec{0, -1},
		}[inst]
		assert(ok, "Unexpected input")

		for i := 0; i < n; i++ {
			rope1.move(dv)
			rope2.move(dv)
		}
	}

	fmt.Printf("Solution Part 1: %d\n", len(rope1.seen))
	fmt.Printf("Solution Part 2: %d\n", len(rope2.seen))
}
