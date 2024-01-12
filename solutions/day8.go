package solutions

import (
	"bufio"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"
)

func Day8(file *os.File) {
	var steps int

	parser := newWastelandParser(file)
	instructions, directions := parser.Parse()
	desert := newDesert(instructions, directions)
	steps = desert.TraverseDesert("AAA", "ZZZ")

	fmt.Printf("Number of steps taken: %v\n", steps)
}

type Desert struct {
	instructions *DirectionRing
	directions MapInstructions
}

func newDesert(i *DirectionRing, d MapInstructions) *Desert {
	desert := &Desert{
		instructions: i,
		directions: d,
	}
	return desert
}

func (d *Desert) TraverseDesert(start, end string) int {
	var count int
	curr := start
	nextDir := d.instructions.head

	slog.Debug("Endpoint", "end", d.directions[end])

	for curr != end{
		slog.Debug("Current Step", "current", curr, "choices", d.directions[curr], "direction", nextDir.literal, "end", end)
		choices, ok := d.directions[curr]
		if !ok {
			log.Fatalf("Direction without mapping")
		}
		curr = choices[nextDir.direction]
		nextDir = nextDir.next
		count++
	}

	return count
}

type DirectionRing struct {
	head   *dRingNode
	tail   *dRingNode
	length int
}

type MapInstructions map[string][]string

func (r DirectionRing) String() string {
	var sb strings.Builder

	curr := r.head
	for i := 0; i < r.length-1; i++ {
		sb.WriteString(curr.literal)
		sb.WriteString("->")
		curr = curr.next
	}
	sb.WriteString(r.tail.literal)

	return sb.String()
}

func (r *DirectionRing) Insert(direction byte) {
	newNode := &dRingNode{}
	newNode.literal = string(direction)
	r.length++

	switch direction {
	case 'R':
		newNode.direction = 1
	case 'L':
		newNode.direction = 0
	}

	if r.head == nil {
		r.head = newNode
		r.head.next = newNode
		r.tail = newNode
		r.tail.next = r.head
		return
	}

	if r.tail != nil {
		r.tail.next = newNode
	}
	r.tail = newNode
	r.tail.next = r.head
	return
}

type dRingNode struct {
	direction int
	literal   string
	next      *dRingNode
}

type WastelandParser struct {
	scanner *bufio.Scanner
}

func newWastelandParser(file *os.File) *WastelandParser {
	p := &WastelandParser{bufio.NewScanner(file)}
	return p
}

func (p *WastelandParser) Parse() (*DirectionRing, MapInstructions) {
	p.scanner.Scan()
	ringLine := p.scanner.Text()
	ring := p.parseRing(ringLine)
	directions := p.parseDirections()
	return ring, directions
}

func (p *WastelandParser) parseRing(line string) *DirectionRing {
	ring := &DirectionRing{}
	for _, b := range []byte(line) {
		ring.Insert(b)
	}
	slog.Debug("Built this ring", "ring", ring)
	return ring
}

func (p *WastelandParser) parseDirections() MapInstructions {
	mapDirections := make(MapInstructions)
	for p.scanner.Scan() {
		line := p.scanner.Text()
		if line == "" {
			continue
		}
		split := strings.Split(line, " = ")
		if len(split) != 2 {
			log.Fatalf("Line not properly formatted: %v\n", line)
		}
		directions := strings.Trim(split[1], "()")
		mapDirections[split[0]] = strings.Split(directions, ", ")
	}
	slog.Debug("Found these directions", "directions", mapDirections)
	return mapDirections
}
