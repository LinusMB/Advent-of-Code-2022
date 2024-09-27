package main

import (
	"fmt"
	"hash/fnv"
	"log"
	"os"
)

func die(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %v", msg, err)
	}
}

type Queue[T any] struct {
	elements []T
}

func (q *Queue[T]) Empty() bool {
	return len(q.elements) == 0
}

func (q *Queue[T]) Enqueue(element T) {
	q.elements = append(q.elements, element)
}

func (q *Queue[T]) Dequeue() (T, bool) {
	if q.Empty() {
		var zero T
		return zero, false
	}
	element := q.elements[0]
	q.elements = q.elements[1:]
	return element, true
}

type Point struct{ x, y int }

func (p Point) Add(other Point) Point {
	return Point{p.x + other.x, p.y + other.y}
}

func flood(grid [][]byte) int {
	miny := len(grid)

	q := Queue[Point]{}

	visited := make(map[Point]bool)

	q.Enqueue(Point{0, len(grid) - 1})

	for !q.Empty() {
		cur, _ := q.Dequeue()

		if cur.x < 0 || cur.x > 6 || cur.y < 0 || cur.y >= len(grid) ||
			grid[cur.y][cur.x] == '#' ||
			visited[cur] {
			continue
		}

		visited[cur] = true

		miny = min(miny, cur.y)

		for _, dir := range []Point{{0, 1}, {1, 0}, {0, -1}, {-1, 0}} {
			q.Enqueue(cur.Add(dir))
		}
	}
	return miny
}

type Block [][]byte

var block1 = Block{
	{'.', '.', '.', '.', '.', '.', '.'},
	{'.', '.', '.', '.', '.', '.', '.'},
	{'.', '.', '.', '.', '.', '.', '.'},
	{'.', '.', '#', '#', '#', '#', '.'},
}

var block2 = Block{
	{'.', '.', '.', '.', '.', '.', '.'},
	{'.', '.', '.', '#', '.', '.', '.'},
	{'.', '.', '#', '#', '#', '.', '.'},
	{'.', '.', '.', '#', '.', '.', '.'},
}

var block3 = Block{
	{'.', '.', '.', '.', '.', '.', '.'},
	{'.', '.', '.', '.', '#', '.', '.'},
	{'.', '.', '.', '.', '#', '.', '.'},
	{'.', '.', '#', '#', '#', '.', '.'},
}

var block4 = Block{
	{'.', '.', '#', '.', '.', '.', '.'},
	{'.', '.', '#', '.', '.', '.', '.'},
	{'.', '.', '#', '.', '.', '.', '.'},
	{'.', '.', '#', '.', '.', '.', '.'},
}

var block5 = Block{
	{'.', '.', '.', '.', '.', '.', '.'},
	{'.', '.', '.', '.', '.', '.', '.'},
	{'.', '.', '#', '#', '.', '.', '.'},
	{'.', '.', '#', '#', '.', '.', '.'},
}

func (b Block) left() func() {
	canmove := true

	for y := 0; y < 4; y++ {
		canmove = canmove && b[y][0] == '.'
	}

	if canmove {
		for y := 0; y < 4; y++ {
			for x := 0; x < 6; x++ {
				b[y][x] = b[y][x+1]
			}
			b[y][6] = '.'
		}
	}
	return func() {
		if canmove {
			b.right()
		}
	}
}

func (b Block) right() func() {
	canmove := true

	for y := 0; y < 4; y++ {
		canmove = canmove && b[y][6] == '.'
	}

	if canmove {
		for y := 0; y < 4; y++ {
			for x := 6; x > 0; x-- {
				b[y][x] = b[y][x-1]
			}
			b[y][0] = '.'
		}
	}
	return func() {
		if canmove {
			b.left()
		}
	}
}

func (b Block) reset() {
	// all blocks start with an x-offset of 2
	for i := 0; i < 6; i++ {
		b.left()
	}
	b.right()
	b.right()
}

type Rotate[T any] struct {
	data []T
	i    int
}

func (r *Rotate[T]) next() T {
	v := r.data[r.i]
	r.i = (r.i + 1) % len(r.data)
	return v
}

type Grid struct {
	grid   [][]byte
	bottom int
	top    int
}

func hashargs(grid [][]byte, patternIndex, blockIndex int) uint64 {
	h := fnv.New64a()

	for i := range grid {
		h.Write(grid[i])
	}

	h.Write([]byte(fmt.Sprintf("%d", patternIndex)))

	h.Write([]byte(fmt.Sprintf("%d", blockIndex)))

	return h.Sum64()
}

func main() {
	var jets Rotate[byte]

	var err error
	jets.data, err = os.ReadFile("input.txt")
	die(err, "Read file")
	jets.data = jets.data[:len(jets.data)-1]
	jets.i = 0

	var blocks Rotate[Block]
	blocks.data = []Block{block1, block2, block3, block4, block5}
	blocks.i = 0

	g := Grid{[][]byte{{'.', '.', '.', '.', '.', '.', '.'}}, 0, 0}

	cache := make(map[uint64]struct{ top, when int })

	now := 0

	simulate := func(until int) {
		for now < until {
			block := blocks.next()

			step(&g, block, &jets)

			hash := hashargs(g.grid, jets.i, blocks.i)

			if v, ok := cache[hash]; ok {
				growth := g.top - v.top

				timediff := now - v.when

				timeleft := until - now

				g.top += (timeleft / timediff) * growth
				g.bottom += (timeleft / timediff) * growth

				now += (timeleft / timediff) * timediff
			} else {
				cache[hash] = struct{ top, when int }{g.top, now}
			}
			block.reset()
			now += 1
		}
	}

	simulate(2022)

	fmt.Printf("Solution Part 1: %v\n", g.top)

	simulate(1000000000000)

	fmt.Printf("Solution Part 2: %v\n", g.top)
}

func step(g *Grid, block Block, jets *Rotate[byte]) {
	for i := 0; i < 8; i++ {
		g.grid = append(g.grid, []byte{'.', '.', '.', '.', '.', '.', '.'})
	}

	ycur := g.top + 3

	collision := func() bool {
		for y := 0; y < 4; y++ {
			for x := 0; x < 7; x++ {
				if ycur < g.bottom ||
					(g.grid[ycur-g.bottom+y][x] == '#' && block[3-y][x] == '#') {
					return true
				}
			}
		}
		return false
	}

	update := func() {
		for y := 0; y < 4; y++ {
			for x := 0; x < 7; x++ {
				if block[3-y][x] == '#' {
					g.grid[ycur-g.bottom+y][x] = '#'
					g.top = max(g.top, ycur+y+1)
				}
			}
		}
		minreach := flood(g.grid)
		g.grid = g.grid[minreach:(g.top - g.bottom)]
		g.bottom += minreach
	}

	var undo func()
	for {
		dir := jets.next()
		if dir == '<' {
			undo = block.left()
		} else {
			undo = block.right()
		}
		if collision() {
			undo()
		}

		ycur -= 1

		if collision() {
			ycur += 1
			update()
			return
		}
	}
}
