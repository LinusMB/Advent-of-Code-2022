package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func die(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %v", msg, err)
	}
}

type Monkey struct {
	id        int
	items     []uint64
	operation struct {
		name    string
		operand string
	}
	divisor   uint64
	toId      map[bool]int
	inspected int
}

func round(monkeys []Monkey, divide bool) {
	var z uint64 = 1

	for _, m := range monkeys {
		z *= m.divisor
	}

	for j, m := range monkeys {
		for _, item := range m.items {

			item %= z

			var op func(uint64, uint64) uint64

			switch m.operation.name {
			case "+":
				op = func(a, b uint64) uint64 {
					return a + b
				}
			case "*":
				op = func(a, b uint64) uint64 {
					return a * b
				}
			}

			if m.operation.operand == "old" {
				item = op(item, item)
			} else {
				opd, err := strconv.ParseUint(m.operation.operand, 10, 64)
				die(err, "round")
				item = op(item, opd)
			}

			if divide {
				item /= 3
			}

			toId := m.toId[item%m.divisor == 0]
			monkeys[toId].items = append(monkeys[toId].items, item)

			monkeys[j].inspected += 1
		}
		monkeys[j].items = monkeys[j].items[:0]
	}
}

func business(monkeys []Monkey) int {
	top1, top2 := 0, 0
	for _, m := range monkeys {
		if m.inspected > top1 {
			top2 = top1
			top1 = m.inspected
		} else if m.inspected > top2 {
			top2 = m.inspected
		}
	}
	return top1 * top2
}

func main() {
	file, err := os.Open("input.txt")
	die(err, "Open file")
	defer file.Close()

	sc := bufio.NewScanner(file)

	var monkeys1 []Monkey
	var monkeys2 []Monkey

	for sc.Scan() {
		if len(sc.Text()) == 0 {
			continue
		}

		var (
			line string
			m    Monkey
		)

		line = sc.Text()
		fmt.Sscanf(line, "Monkey %d:", &m.id)

		sc.Scan()
		line = sc.Text()

		re := regexp.MustCompile(`\d+`)

		matches := re.FindAllString(line, -1)

		for i := range matches {
			num, err := strconv.ParseUint(matches[i], 10, 64)
			die(err, "sc.Scan")
			m.items = append(m.items, num)
		}

		sc.Scan()
		line = strings.TrimSpace(sc.Text())
		fmt.Sscanf(line, "Operation: new = old %s %s", &m.operation.name, &m.operation.operand)

		sc.Scan()
		line = strings.TrimSpace(sc.Text())
		fmt.Sscanf(line, "Test: divisible by %d:", &m.divisor)

		sc.Scan()

		m.toId = make(map[bool]int)

		var id int
		line = strings.TrimSpace(sc.Text())
		fmt.Sscanf(line, "If true: throw to monkey %d", &id)
		m.toId[true] = id

		sc.Scan()
		line = strings.TrimSpace(sc.Text())
		fmt.Sscanf(line, "If false: throw to monkey %d", &id)
		m.toId[false] = id

		monkeys1 = append(monkeys1, m)

		var m2 Monkey
		m2.id = m.id
		m2.items = make([]uint64, len(m.items))
		copy(m2.items, m.items)
		m2.operation = m.operation
		m2.divisor = m.divisor
		m2.toId = m.toId
		m2.inspected = m.inspected

		monkeys2 = append(monkeys2, m2)
	}

	for i := 0; i < 20; i++ {
		round(monkeys1, true)
	}

	for i := 0; i < 10000; i++ {
		round(monkeys2, false)
	}

	fmt.Printf("Solution Part 1: %v\n", business(monkeys1))
	fmt.Printf("Solution Part 2: %v\n", business(monkeys2))
}
