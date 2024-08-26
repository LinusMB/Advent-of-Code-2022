package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
)

func die(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %v", msg, err)
	}
}

type F struct {
	name    string
	size    int
	isDir   bool
	content map[string]*F
}

func (f *F) AddSize(s int) {
	f.size += s
	if p, ok := f.content[".."]; ok {
		p.AddSize(s)
	}
}

func main() {
	file, err := os.Open("input.txt")
	die(err, "Open file")
	defer file.Close()

	sc := bufio.NewScanner(file)

	root := &F{
		name:    "/",
		isDir:   true,
		content: map[string]*F{},
	}

	cur := &F{
		content: map[string]*F{
			"/": root,
		},
	}

	dirs := make(map[*F]bool)

	for sc.Scan() {
		line := sc.Text()

		switch {
		case strings.HasPrefix(line, "$ cd"):
			var name string
			fmt.Sscanf(line, "$ cd %s", &name)
			cur = cur.content[name]
			dirs[cur] = true
		case strings.HasPrefix(line, "$ ls"):
			continue
		case strings.HasPrefix(line, "dir"):
			var name string
			fmt.Sscanf(line, "dir %s", &name)
			f := &F{name: name, isDir: true, content: map[string]*F{"..": cur}}
			cur.content[name] = f
		default:
			var (
				size int
				name string
			)
			fmt.Sscanf(line, "%d %s", &size, &name)
			f := &F{name: name, size: size, content: map[string]*F{"..": cur}}
			cur.content[name] = f
			cur.AddSize(size)
		}
	}

	need := 30000000 - (70000000 - root.size)

	sol1, sol2 := 0, math.MaxInt

	for f := range dirs {
		if f.size > need && f.size < sol2 {
			sol2 = f.size
		}
		if f.size < 100000 {
			sol1 += f.size
		}
	}

	fmt.Printf("Solution Part 1: %d\n", sol1)
	fmt.Printf("Solution Part 2: %d\n", sol2)
}
