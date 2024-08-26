package main

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"os"

	"golang.org/x/exp/constraints"
)

func die(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %v", msg, err)
	}
}

type Coordinate struct {
	nRows, nCols, row, col int
}

func newCoordinateFromIndex(nRows, nCols, i int) Coordinate {
	return Coordinate{
		nRows: nRows,
		nCols: nCols,
		row:   i / nCols,
		col:   i % nCols,
	}
}

func (c Coordinate) index() int {
	return c.row*c.nCols + c.col
}

func (c Coordinate) add(row, col) Coordinate {
	return Coordinate{
		nRows: c.nRows,
		nCols: c.nCols,
		row:   c.row + row,
		col:   c.col + col,
	}
}

func (c Coordinate) candidates() []Coordinate {
	var cs []Coordinate

	if c.col > 0 {
		cs = append(cs, c.add(0, -1))
	}
	if c.col < c.nCols-1 {
		cs = append(cs, c.add(0, 1))
	}
	if c.row > 0 {
		cs = append(cs, c.add(-1, 0))
	}
	if c.row < c.nRows-1 {
		cs = append(cs, c.add(1, 0))
	}
	return cs
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

func min[T constraints.Ordered](x, y T) T {
	if x < y {
		return x
	}
	return y
}

func shortest(grid []byte, start, end Coordinate) int {
	q := &Queue[Coordinate]{}
	q.Enqueue(start)

	distance := make(map[Coordinate]int)
	distance[start] = 0

	for !q.Empty() {
		last, _ := q.Dequeue()

		if last == end {
			return distance[end]
		}

		for _, c := range last.candidates() {
			if _, ok := distance[c]; ok || int(grid[c.index()])-int(grid[last.index()]) > 1 {
				continue
			}
			distance[c] = distance[last] + 1

			q.Enqueue(c)
		}
	}
	return math.MaxInt
}

func main() {
	data, err := os.ReadFile("input.txt")
	die(err, "Read File")

	nRows := bytes.Count(data, []byte("\n"))
	nCols := bytes.IndexByte(data, '\n')

	grid := bytes.ReplaceAll(data, []byte("\n"), []byte(""))

	var (
		start Coordinate
		end   Coordinate
	)
	{
		var i int
		i = bytes.IndexByte(grid, 'S')
		start = newCoordinateFromIndex(nRows, nCols, i)
		grid[i] = 'a'
		i = bytes.IndexByte(grid, 'E')
		grid[i] = 'z'
		end = newCoordinateFromIndex(nRows, nCols, i)
	}

	fmt.Printf("Solution Part 1: %v\n", shortest(grid, start, end))

	minDist := math.MaxInt

	for i := 0; i < len(grid); i++ {
		if grid[i] == 'a' {
			minDist = min(minDist, shortest(grid, newCoordinateFromIndex(nRows, nCols, i), end))
		}
	}

	fmt.Printf("Solution Part 2: %v\n", minDist)
}
