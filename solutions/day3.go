package solutions

import (
	"bufio"
	"fmt"
	"log/slog"
	"io"
	"os"

	"github.com/ryanpdenoux/advent-of-code/utils"
)

func Day3(file *os.File) {
	var (
		parts   int = 0
		gears   int = 0
		gearAcc int = 1
	)

	schematic := newSchematicFromFile(file)
	numbers := schematic.findPartNumbers()
	for _, number := range numbers {
		parts += number
	}

	fmt.Printf("The sum of the engine parts is %d\n", parts)

	for _, v := range schematic.gears {
		if len(v) <= 1 {
			continue
		}

		for _, e := range v {
			gearAcc *= e
		}
		gears += gearAcc
		gearAcc = 1
	}

	fmt.Printf("The sum of gear ratios is %d\n", gears)
}

type EngineSchematic struct {
	reader *bufio.Reader
	row    int
	gears  map[Part][]int
	prev   []byte
	curr   []byte
	next   []byte
}

type Part struct {
	point Point
	gear bool
}


type Point struct {
	char byte
	x int
	y int
}

func (p Point) String() string {
	return fmt.Sprintf("char: %v (%3d, %3d)", string(p.char), p.x, p.y)
}

// Creates an instance from a file handle and sets up the current and next lines
func newSchematicFromFile(file *os.File) *EngineSchematic {
	s := &EngineSchematic{}
	s.reader = bufio.NewReader(file)
	s.row = 1
	s.gears = make(map[Part][]int)
	s.advanceLines()

	return s
}

func (s *EngineSchematic) advanceLines() bool {
	s.row++
	curr := s.curr
	next := s.next
	s.curr = next
	s.prev = curr

	next, err := s.reader.ReadBytes('\n')
	if err == io.EOF {
		s.next = []byte{}
		return true
	}

	s.next = next[:len(next)-1]
	return true
}

func (s *EngineSchematic) findPartNumbers() []int {
	partNumbers := []int{}

	// Here we must check if there is a current line
	for s.advanceLines() && len(s.curr) > 0 {
		rowNumbers := s.findRowNumbers()
		partNumbers = append(partNumbers, rowNumbers...)
	}

	return partNumbers
}

func (s *EngineSchematic) findRowNumbers() []int {
	numbers := []int{}

	for i := 0; i < len(s.curr); i++ {
		number, ok := utils.FindNumberInBytes(s.curr, i)
		if ok {
			part, ok := s.checkNumber(number, i)
			if ok {
				numbers = append(numbers, number)
				if part.gear {
					s.gears[part] = append(s.gears[part], number)
				}
			}
			// advance index by length of digits
			i += utils.LengthOfInt(number) - 1
		}
	}

	slog.Debug("Found these part numbers", "partNumbers", numbers)
	return numbers
}

// Check area surrounding number
func (s *EngineSchematic) checkNumber(number, pos int) (Part, bool) {
	numberWidth := utils.LengthOfInt(number)
	points := s.getPointsToCheck(pos, numberWidth)
	slog.Debug("Checking points for symbols", "points", points)

	return checkPoints(points)
}

func (s *EngineSchematic) getPointsToCheck(pos, width int) []Point {
	points := []Point{}

	for _, i := range utils.MakeRange(pos-1, pos+width+1) {
		if i >= 0 && i < len(s.curr) {
			point := Point{}
			point.y = i

			// first line if false
			if len(s.prev) > 0 {
				point.char = s.prev[i]
				point.x = s.row - 1
				points = append(points, point)
			}

			point.char = s.curr[i]
			point.x = s.row
			points = append(points, point)

			// last line if false
			if len(s.next) > 0 {
				point.char = s.next[i]
				point.x = s.row + 1
				points = append(points, point)
			}
		}
	}

	return points
}

func checkPoints(points []Point) (Part, bool) {
	for _, point := range points {
		if !utils.IsDigit(point.char) && point.char != '.' {
			part := Part{point: point}
			if point.char == '*' {
				part.gear = true
			}
			return part, true
		}
	}
	return Part{}, false
}
