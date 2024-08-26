package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
)

func die(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %v", msg, err)
	}
}

type Packet struct {
	id      int
	integer int
	isList  bool
	list    []Packet
}

func parsePacket(data string) (string, Packet) {
	var p Packet

	for {
		if data == "" {
			return "", p
		}
		if data[0] != ',' {
			break
		}
		data = data[1:]
	}

	if data[0] == '[' {
		p.isList = true
		data = data[1:]
		for data[0] != ']' {
			var p0 Packet
			data, p0 = parsePacket(data)
			p.list = slices.Insert(p.list, len(p.list), p0)
		}
		return data[1:], p
	}

	i := 0
	for ; data[i] >= '0' && data[i] <= '9'; i++ {
	}

	var err error
	p.isList = false
	p.integer, err = strconv.Atoi(data[:i])
	die(err, "parse packet")

	return data[i:], p
}

type Order int

const (
	lt = -1
	gt = 1
	eq = 0
)

func getOrder(left, right Packet) Order {
	if !left.isList && !right.isList {
		if left.integer < right.integer {
			return lt
		}
		if left.integer > right.integer {
			return gt
		}
		return eq
	}

	if !left.isList {
		old := left
		left.list = slices.Insert(left.list, len(left.list), old)
		left.isList = true
	}

	if !right.isList {
		old := right
		right.list = slices.Insert(right.list, len(right.list), old)
		right.isList = true
	}

	if left.isList && right.isList {
		for i := 0; i < len(left.list) && i < len(right.list); i++ {
			o := getOrder(left.list[i], right.list[i])
			if o != eq {
				return o
			}
		}
		if len(left.list) < len(right.list) {
			return lt
		}
		if len(left.list) > len(right.list) {
			return gt
		}
		return eq
	}
	panic("Unreachable")
}

func main() {
	file, err := os.Open("input.txt")
	die(err, "Open file")
	defer file.Close()

	sc := bufio.NewScanner(file)

	id := 0
	var packets []Packet

	for sc.Scan() {
		if len(sc.Text()) == 0 {
			continue
		}

		_, p := parsePacket(sc.Text())
		p.id = id
		id += 1
		packets = slices.Insert(packets, len(packets), p)
	}

	sum := 0
	for i := 0; i < len(packets); i += 2 {
		if getOrder(packets[i], packets[i+1]) == lt {
			sum += (i / 2) + 1
		}
	}

	_, divider1 := parsePacket("[[2]]")
	_, divider2 := parsePacket("[[6]]")

	divider1.id = id
	id += 1
	divider2.id = id
	id += 1

	packets = slices.Insert(packets, len(packets), divider1)
	packets = slices.Insert(packets, len(packets), divider2)

	slices.SortFunc(packets, func(a, b Packet) int {
		return int(getOrder(a, b))
	})

	prod := 1
	for i := 0; i < len(packets); i++ {
		if packets[i].id == divider1.id || packets[i].id == divider2.id {
			prod *= i + 1
		}
	}

	fmt.Printf("Solution Part 1: %v\n", sum)
	fmt.Printf("Solution Part 2: %v\n", prod)
}
