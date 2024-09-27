package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"slices"

	"golang.org/x/exp/constraints"
)

func die(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %v", msg, err)
	}
}

func abs[T constraints.Signed | constraints.Float](x T) T {
	if x < 0 {
		return -x
	}
	return x
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

type Sensor struct {
	Point
	radius int
	beacon Point
}

func manhattan(a, b Point) int {
	return abs(a.x-b.x) + abs(a.y-b.y)
}

func visible(sensors []Sensor, p Point) bool {
	return residual(sensors, p) >= 0
}

func residual(sensors []Sensor, p Point) int {
	maxv := math.MinInt
	for i := range sensors {
		maxv = max(maxv, sensors[i].radius-manhattan(sensors[i].Point, p))
	}
	return maxv
}

var notfound = Point{-1, -1}

func search(sensors []Sensor, a, b Point) Point {
	if manhattan(a, b) <= 5 {
		for t := a; ; t.x, t.y = t.x+1, t.y+1 {
			if !visible(sensors, t) {
				return t
			}
			if t.Eq(b) {
				return notfound
			}
		}
	}

	if (residual(sensors, a) + residual(sensors, b)) >= manhattan(a, b) {
		return notfound
	}

	mid := Point{(a.x + b.x) / 2, (a.y + b.y) / 2}

	if found := search(sensors, a, mid); !found.Eq(notfound) {
		return found
	}

	return search(sensors, mid, b)
}

func main() {
	file, err := os.Open("input.txt")
	die(err, "Open file")
	defer file.Close()

	sc := bufio.NewScanner(file)

	var sensors []Sensor

	minX := math.MaxInt
	maxX := math.MinInt

	for sc.Scan() {
		line := sc.Text()

		var sensor Sensor
		fmt.Sscanf(line, "Sensor at x=%d, y=%d: closest beacon is at x=%d, y=%d",
			&sensor.x,
			&sensor.y,
			&sensor.beacon.x,
			&sensor.beacon.y)
		sensor.radius = manhattan(sensor.Point, sensor.beacon)
		sensors = slices.Insert(sensors, len(sensors), sensor)

		minX = min(minX, sensor.x-sensor.radius)
		maxX = max(maxX, sensor.x+sensor.radius)
	}

	count := 0
Loop1:
	for x := minX; x <= maxX; x++ {
		for _, s := range sensors {
			if manhattan(s.Point, Point{x, 2000000}) <= s.radius &&
				!s.beacon.Eq(Point{x, 2000000}) {
				count += 1
				continue Loop1
			}
		}
	}

	fmt.Printf("Solution Part 1: %v\n", count)

	var frequency int

Loop2:
	for _, s := range sensors {
		t := Point{s.x, s.y + s.radius + 1}
		ml := Point{s.x - s.radius - 1, s.y}
		mr := Point{s.x + s.radius + 1, s.y}
		b := Point{s.x, s.y - s.radius - 1}

		for _, line := range [][]Point{{ml, t}, {mr, t}, {b, ml}, {b, mr}} {
			if found := search(sensors, line[0], line[1]); !found.Eq(notfound) &&
				found.x >= 0 &&
				found.x <= 4000000 &&
				found.y >= 0 &&
				found.y <= 4000000 {
				frequency = found.x*4000000 + found.y
				break Loop2
			}
		}
	}
	fmt.Printf("Solution Part 2: %v\n", frequency)
}
