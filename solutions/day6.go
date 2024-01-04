package solutions

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/ryanpdenoux/advent-of-code/utils"
)

func Day6(file *os.File) {
	var accumulatedRecord int = 1

	parser := createRegattaParser(file)
	records := parser.Parse()
	boat := &RegattaBoat{1}

	for _, record := range records {
		numBetterRecords := boat.AttemptRace(record)
		accumulatedRecord = accumulatedRecord * numBetterRecords
	}

	fmt.Printf("Accumulated Records Beaten: %v\n", accumulatedRecord)
}

type RegattaBoat struct {
	baseVelocity int
}

// Returns number of ways in which a race can be won by brute force
// Only works if charging is linear
func (b *RegattaBoat) AttemptRace(record RegattaRecord) int {
	var (
		right int
		left int
	)

	for i := record.Time; i != 0; i-- {
		distance := b.velocity(record.Time-i) * i
		if distance > record.Distance {
			left = i
		}
	}

	for i := 0; i < record.Time; i++ {
		distance := b.velocity(record.Time-i) * i
		if distance > record.Distance {
			right = i
		}
	}

	slog.Info("Boat was able to beat record n times", "n", right-left)
	return right-left+1
}

func (b *RegattaBoat) velocity(chargeTime int) int {
	return b.baseVelocity * chargeTime
}

type RegattaRecord struct {
	Time int
	Distance int
}

type RegattaParser struct {
	scanner *bufio.Scanner
	rawTime string
	rawDist string
}

func createRegattaParser(file io.Reader) *RegattaParser {
	p := &RegattaParser{}
	p.scanner = bufio.NewScanner(file)
	return p
}

func (p *RegattaParser) Parse() []RegattaRecord {
	records := []RegattaRecord{}

	strTimes := p.prepareLine()
	strDists := p.prepareLine()
	times := utils.StringSliceToIntSlice(strTimes)
	dists := utils.StringSliceToIntSlice(strDists)
	slog.Debug("Found times and dists", "times", times, "dists", dists)

	for i := 0; i < len(times); i++ {
		records = append(records, RegattaRecord{times[i], dists[i]})
	}

	return records
}

func (p *RegattaParser) AlternateParse() RegattaRecord {
	record := RegattaRecord{}

	strTimes := p.prepareLine()
	strTime := strings.Join(strTimes, "")
	slog.Debug("Raw time", "Time", strTime)
	time, err := strconv.Atoi(strTime)
	if err != nil {
		log.Fatalf("Not able to convert fields %v: %v\n", strTimes, err)
	}
	record.Time = time
	slog.Debug("Found time", "time", time)

	strDists := p.prepareLine()
	strDist := strings.Join(strDists, "")
	dist, err := strconv.Atoi(strDist)
	if err != nil {
		log.Fatalf("Not able to convert fields %v: %v\n", strDists, err)
	}
	record.Distance = dist

	return record
}

func (p *RegattaParser) prepareLine() []string {
	p.scanner.Scan()
	rawLine := p.scanner.Text()
	line := strings.Split(rawLine, ":")
	return strings.Fields(line[1])
}
