package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strings"

	"golang.org/x/exp/constraints"
)

func die(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %v", msg, err)
	}
}

func max[T constraints.Ordered](x, y T) T {
	if x > y {
		return x
	}
	return y
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

func main() {
	file, err := os.Open("input.txt")
	die(err, "Open file")
	defer file.Close()

	var lines []string
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}

	strids := make(map[string]int, len(lines))
	rates := make([]int, len(lines))
	graph := make([][]int, len(lines))
	for i := range graph {
		graph[i] = make([]int, len(lines))
	}

	for i := range lines {
		var (
			strid string
			rate  int
		)
		fmt.Sscanf(lines[i], "Valve %s has flow rate=%d", &strid, &rate)

		strids[strid] = i
		rates[i] = rate
	}

	for i := range lines {
		index := strings.Index(lines[i], ";")

		re := regexp.MustCompile(`[A-Z]{2}`)

		matches := re.FindAllString(lines[i][index+1:], -1)

		for j := range matches {
			graph[i][strids[matches[j]]] = 1
		}
	}

	dist := mindist(graph)
	state := make([]bool, len(lines))
	cache := make(map[uint64]int)

	search(dist, rates, state, strids["AA"], 30, 0, cache)

	pressure1 := 0

	for state := range cache {
		pressure1 = max(pressure1, cache[state])
	}

	fmt.Printf("Solution Part 1: %v\n", pressure1)

	cache = make(map[uint64]int)

	search(dist, rates, state, strids["AA"], 26, 0, cache)

	pressure2 := 0

	for state1 := range cache {
		for state2 := range cache {
			if (state1 & state2) == 0 {
				pressure2 = max(pressure2, cache[state1]+cache[state2])
			}
		}
	}

	fmt.Printf("Solution Part 2: %v\n", pressure2)
}

func bfs(graph [][]int, start int) []int {
	dist := make([]int, len(graph))
	for i := range dist {
		dist[i] = math.MaxInt
	}
	dist[start] = 0

	q := Queue[int]{}
	q.Enqueue(start)

	for !q.Empty() {
		cur, _ := q.Dequeue()

		for node := 0; node < len(graph); node++ {
			if graph[cur][node] == 1 && dist[node] == math.MaxInt {
				dist[node] = dist[cur] + 1
				q.Enqueue(node)
			}
		}
	}
	return dist
}

func mindist(graph [][]int) [][]int {
	dist := make([][]int, len(graph))

	for node := 0; node < len(graph); node++ {
		dist[node] = bfs(graph, node)
	}

	return dist
}

func tobitmask(state []bool) uint64 {
	var bitmask uint64

	for i := range state {
		if state[i] {
			bitmask |= 1 << i
		}
	}
	return bitmask
}

func search(dist [][]int, rates []int, state []bool, cur int, minutes int, pressure int, cache map[uint64]int) {
	bitmask := tobitmask(state)

	cache[bitmask] = max(cache[bitmask], pressure)

	for node := 0; node < len(state); node++ {
		remaining := minutes - dist[cur][node] - 1

		if rates[node] <= 0 || state[node] || remaining < 0 {
			continue
		}

		state[node] = true

		search(dist, rates, state, node, remaining, pressure+rates[node]*remaining, cache)

		state[node] = false
	}
}
