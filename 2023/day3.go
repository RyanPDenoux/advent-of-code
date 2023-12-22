package solution

import (
	"bufio"
	"fmt"
	"math"
	"os"
	// "slices"
	"strconv"
)

func Day3(file *os.File) {
	var sum int

	// fmt.Printf("This is my file: %v\n", file)
	schematic := newSchematicFromFile(file)
	numbers := schematic.FindAllNumbers()
	for _, number := range numbers {
		sum += number
	}

	fmt.Printf("The sum of the engine parts is %d\n", sum)
}

type EngineSchematic struct {
	scanner *bufio.Scanner
	prev    []byte
	curr    []byte
	next    []byte
}

// Creates an instance from a file handle and sets up the current and next lines
func newSchematicFromFile(file *os.File) *EngineSchematic {
	s := &EngineSchematic{}
	s.scanner = bufio.NewScanner(file)
	s.advanceLines()

	return s
}

func (s *EngineSchematic) advanceLines() bool {
	curr := s.curr
	next := s.next
	s.curr = next
	s.prev = curr

	if s.scanner.Scan() {
		s.next = s.scanner.Bytes()
		return true
	} else {
		// fmt.Printf("Last line: %v\n", s.curr)
		s.next = []byte{}
		return true
	}
}

func (s *EngineSchematic) FindAllNumbers() []int {
	numbers := []int{}

	for s.advanceLines() && len(s.curr) > 0 {
		rowNumbers := s.findNumbers()
		numbers = append(numbers, rowNumbers...)
		// for _, number := range(rowNumbers) {
		// 	if !slices.Contains(numbers, number) {
		// 		numbers = append(numbers, number)
		// 	}
		// }
	}

	return numbers
}

func (s *EngineSchematic) findNumbers() []int {
	numbers := []int{}

	for i := 0; i < len(s.curr); i++ {
		number, ok := determineNumber(s.curr, i)
		if ok {
			if s.isEnginePart(number, i) {
				numbers = append(numbers, number)
			}
			// advance counter by length of digits
			i += lengthOfInt(number)
		}
	}

	fmt.Printf("Found these numbers: %v\n", numbers)
	return numbers
}

func (s *EngineSchematic) isEnginePart(number, pos int) bool {
	points := makeRange(pos-1, pos+lengthOfInt(number)+1)
	validPoints := s.filterPoints(points)

	if len(s.prev) > 0 && len(validPoints) > 0 {
		if containsSymbol(s.prev, validPoints) {
			return true
		}
	}
	if len(s.next) > 0 && len(validPoints) > 0 {
		if containsSymbol(s.next, validPoints) {
			return true
		}
	}
	if containsSymbol(s.curr, validPoints) {
		return true
	}

	return false
}

func (s *EngineSchematic) filterPoints(points []int) []int {
	filtered := []int{}

	for _, point := range points {
		if point > 0 && point < len(s.curr) {
			filtered = append(filtered, point)
		}
	}

	// fmt.Printf("Points before: %v\nPoints after: %v\n", points, filtered)
	return filtered
}

func determineNumber(data []byte, pos int) (int, bool) {
	var num []byte

	// check length first to avoid index errors
	for pos < len(data) && isDigit(data[pos]) {
		num = append(num, data[pos])
		pos++
	}

	intNum, err := strconv.Atoi(string(num))
	if err != nil {
		return 0, false
	}

	return intNum, true
}

func isSymbol(char byte) bool {
	return !isDigit(char) && char != '.'
	// return char != '.'
}

func lengthOfInt(number int) int {
	return int(math.Log10(float64(number))) + 1
}

func makeRange(min, max int) []int {
	a := make([]int, max-min)
	for i := range a {
		a[i] = min + i
	}
	return a
}

func containsSymbol(data []byte, points []int) bool {
	// fmt.Printf("Check this data:%v\nWith these points:%v", data, points)
	for _, point := range points {
		if isSymbol(data[point]) {
			return true
		}
	}
	return false
}

// func dedupSlice(slice []int) []int {
// 	exists := []int{}
// 	deduped := []int{}

// 	for _, number := range(slice) {
// 		if !slices.Contains(exists, number) {
// 			deduped = append(deduped, number)
// 			exists = append(exists, number)
// 		}
// 	}
// 	return deduped
// }
