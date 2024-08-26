package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func die(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %v", msg, err)
	}
}

type Crate struct {
	id   byte
	next *Crate
	prev *Crate
}

func (c *Crate) Tail() *Crate {
	cur := c
	for cur.next != nil {
		cur = cur.next
	}
	return cur
}

func (c *Crate) Last(n int) *Crate {
	cur := c
	for ; cur.prev != nil && n > 1; n-- {
		cur = cur.prev
	}
	return cur
}

func (c *Crate) Insert(id byte) {
	newCrate := &Crate{id: id}
	if c.next != nil {
		c.next.prev = newCrate
	}
	newCrate.next = c.next
	newCrate.prev = c
	c.next = newCrate
}

func (c *Crate) Move(to *Crate) {
	to.next = c
	c.prev.next = nil
	c.prev = to
}

func stringTails(heads []Crate) string {
	var b strings.Builder
	for _, h := range heads {
		b.WriteByte(h.Tail().id)
	}
	return b.String()
}

func main() {
	file, err := os.Open("input.txt")
	die(err, "Open file")
	defer file.Close()

	sc := bufio.NewScanner(file)

	var (
		heads1 []Crate
		heads2 []Crate
	)

	for sc.Scan() {
		line := sc.Text()

		switch _line := strings.TrimSpace(line); {
		case strings.HasPrefix(_line, "["):
			for i, n := 1, 0; i < len(line); i, n = i+4, n+1 {
				if len(heads1) <= n {
					heads1 = append(heads1, Crate{id: '_'})
					heads2 = append(heads2, Crate{id: '_'})
				}
				if line[i] != ' ' {
					heads1[n].Insert(line[i])
					heads2[n].Insert(line[i])
				}
			}
		case strings.HasPrefix(_line, "move"):
			var amt, src, dst int
			fmt.Sscanf(line, "move %d from %d to %d", &amt, &src, &dst)

			for i := 0; i < amt; i++ {
				crateSrc := heads1[src-1].Tail()
				crateDst := heads1[dst-1].Tail()
				crateSrc.Move(crateDst)
			}
			{
				crateSrc := heads2[src-1].Tail().Last(amt)
				crateDst := heads2[dst-1].Tail()
				crateSrc.Move(crateDst)
			}

		default:
			continue
		}
	}

	fmt.Printf("Solution Part 1: %s\n", stringTails(heads1))
	fmt.Printf("Solution Part 2: %s\n", stringTails(heads2))
}
