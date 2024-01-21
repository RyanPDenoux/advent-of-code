package solutions

import (
	"bufio"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strconv"
	"strings"
)

func Day9(file *os.File) {
	var sum int

	parser := newOasisParser(file)
	records := parser.Parse()
	for _, record := range records {
		sum += record.Predict()
	}

	fmt.Printf("Sum of predicted values: %v\n", sum)
}

type OasisRecord []int

func (r *OasisRecord) Predict() int {
	var prediction int

	prediction = r.predictValue()
	slog.Debug("Predicted value for current record", "record", r, "prediction", prediction)

	return prediction
}

func (r OasisRecord) predictValue() int {
	new := OasisRecord{}
	curr := r[0]
	zero := true

	for i := 1; i < len(r); i++ {
		next := r[i]
		val := next - curr
		new = append(new, val)
		curr = next
		if curr != 0 {
			zero = false
		}
	}

	if zero {
		return 0
	}

	return curr + new.predictValue()
}

// Parsing
type OasisParser struct {
	scanner *bufio.Scanner
}

func newOasisParser(file *os.File) *OasisParser {
	parser := &OasisParser{
		scanner: bufio.NewScanner(file),
	}
	return parser
}

func (p *OasisParser) Parse() []OasisRecord {
	records := []OasisRecord{}

	for p.scanner.Scan() {
		line := p.scanner.Text()
		record := OasisRecord{}
		for _, str := range strings.Fields(line) {
			if val, err := strconv.Atoi(str); err == nil {
				record = append(record, val)
			} else {
				log.Fatalf("Could not parse record: %v\n", line)
			}
		}
		records = append(records, record)
	}

	return records
}
