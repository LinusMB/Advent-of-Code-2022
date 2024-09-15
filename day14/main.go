package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"slices"

	"golang.org/x/exp/constraints"
)

func die(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %v", msg, err)
	}
}

func sign[T constraints.Signed](x T) T {
	if x > 0 {
		return 1
	} else if x < 0 {
		return -1
	}
	return 0
}

func min[T constraints.Ordered](x, y T) T {
	if x < y {
		return x
	}
	return y
}

func max[T constraints.Ordered](x, y T) T {
	if x > y {
		return x
	}
	return y
}

type Point struct{ x, y int }

func (p Point) Eq(other Point) bool {
	return p.x == other.x && p.y == other.y
}

func (p Point) Sub(other Point) Point {
	return Point{p.x - other.x, p.y - other.y}
}

func (p Point) Down() Point {
	return Point{p.x, p.y + 1}
}

func (p Point) DownLeft() Point {
	return Point{p.x - 1, p.y + 1}
}

func (p Point) DownRight() Point {
	return Point{p.x + 1, p.y + 1}
}

type Grid struct {
	grid   []uint8
	source Point
	xsize  int
	ysize  int
}

func (g *Grid) falling(p Point) bool {
	return p.y >= g.ysize-2
}

func (g *Grid) bounds(p Point) bool {
	return p.x >= 0 && p.x < g.xsize && p.y >= 0 && p.y < g.ysize
}

func (g *Grid) free(p Point) bool {
	return g.grid[p.x+p.y*g.xsize] == 0
}

func (g *Grid) step(cur Point) Point {
	candidates := []Point{cur.Down(), cur.DownLeft(), cur.DownRight()}
	for _, cand := range candidates {
		if g.bounds(cand) && g.free(cand) {
			return cand
		}
	}
	return cur
}

func (g *Grid) settle() Point {
	cur := g.source

	for {
		next := g.step(cur)

		if next.Eq(cur) {
			return cur
		}
		cur = next
	}
}

func makegrid(paths [][]Point, source Point) Grid {
	minX, maxX := source.x, source.x
	minY, maxY := source.y, source.y

	for _, path := range paths {
		for _, point := range path {
			minX = min(minX, point.x)
			minY = min(minY, point.y)
			maxX = max(maxX, point.x)
			maxY = max(maxY, point.y)
		}
	}

	maxY += 2

	b := 3 * (maxY - source.y) // upper bound for the base of a triangle formed by the sand

	minX = min(minX, source.x-b/2)
	maxX = max(maxX, source.x+b/2)

	xsize := maxX - minX + 1
	ysize := maxY - minY + 1

	for i := range paths {
		for j := range paths[i] {
			paths[i][j] = paths[i][j].Sub(Point{minX, minY})
		}
	}

	source = source.Sub(Point{minX, minY})

	grid := make([]uint8, xsize*ysize)

	for _, path := range paths {
		for i := 0; i < len(path)-1; i++ {
			pa := path[i]
			pb := path[i+1]

			sx := sign(pb.x - pa.x)
			sy := sign(pb.y - pa.y)

			for x, y := pa.x, pa.y; ; x, y = x+sx, y+sy {
				grid[x+y*xsize] = 1
				if x == pb.x && y == pb.y {
					break
				}
			}
		}
	}

	for x := 0; x < xsize; x++ {
		grid[x+(ysize-1)*xsize] = 1
	}

	return Grid{grid, source, xsize, ysize}
}

func main() {
	file, err := os.Open("input.txt")
	die(err, "Open file")
	defer file.Close()

	sc := bufio.NewScanner(file)

	var paths [][]Point

	for sc.Scan() {
		line := sc.Text()

		re := regexp.MustCompile(`\d+,\d+`)
		matches := re.FindAllString(line, -1)

		points := make([]Point, len(matches))
		for i := range matches {
			fmt.Sscanf(matches[i], "%d,%d", &points[i].x, &points[i].y)
		}
		paths = slices.Insert(paths, len(paths), points)
	}

	count1 := 0

	g := makegrid(paths, Point{500, 0})

	for {
		settled := g.settle()
		if g.falling(settled) {
			break
		}
		g.grid[settled.x+settled.y*g.xsize] = 1
		count1 += 1
	}

	count2 := count1

	for {
		settled := g.settle()
		if settled.Eq(g.source) {
			count2 += 1
			break
		}
		g.grid[settled.x+settled.y*g.xsize] = 1
		count2 += 1
	}

	fmt.Printf("Solution Part 1: %v\n", count1)
	fmt.Printf("Solution Part 2: %v\n", count2)
}
