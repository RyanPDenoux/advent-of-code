package solutions

import (
	"bufio"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"

	"github.com/ryanpdenoux/advent-of-code/utils"
)

func Day5(file *os.File) {
	min := 1024 * 1024 * 1024 * 1024

	parser := createAlmanacParser(file)
	almanac := parser.createAlamanac()
	for _, seed := range almanac.seeds {
		location := almanac.FindLocationToPlant(seed)
		if location < min {
			min = location
		}
	}

	fmt.Printf("Plant here: %v\n", min)

	// location := almanac.FindLocationFromRange()
	// fmt.Printf("Next here: %v\n", location)
}

type Almanac struct {
	seeds                 []int
	seedToSoil            AlmanacMapper
	soilToFertilizer      AlmanacMapper
	fertilizerToWater     AlmanacMapper
	waterToLight          AlmanacMapper
	lightToTemperature    AlmanacMapper
	temperatureToHumidity AlmanacMapper
	humidityToLocation    AlmanacMapper
}

func (a *Almanac) FindLocationToPlant(seed int) int {
	location := seed

	for _, mapper := range a.mappers() {
		location = mapper.findNext(location)
	}
	slog.Info("Mapping Seed to Location", "seed", seed, "location", location)
	return location
}

// func (a *Almanac) FindLocationFromRange() int {
// 	var minimum int

// 	seeds := a.expandSeeds()
// 	slog.Debug("Seeds", "seeds", seeds)
// 	for _, seed := range seeds {
// 		location := seed
// 		for _, mapper := range a.mappers() {
// 			location = mapper.findNext(location)
// 		}
// 		if location < minimum {
// 			minimum = location
// 		}
// 	}

// 	return minimum
// }

// func (a *Almanac) expandSeeds() []int {
// 	seeds := []int{}

// 	for i := 0; i < len(a.seeds) - 1; i = i+2 {
// 		start, length := a.seeds[i], a.seeds[i+1]
// 		for i := start; i < start+length; i++ {
// 			seeds = append(seeds, i)
// 		}
// 	}

// 	return seeds
// }

func (a *Almanac) mappers() []AlmanacMapper {
	return []AlmanacMapper{
		a.seedToSoil,
		a.soilToFertilizer,
		a.fertilizerToWater,
		a.waterToLight,
		a.lightToTemperature,
		a.temperatureToHumidity,
		a.humidityToLocation,
	}
}

type AlmanacMapper struct {
	mappings []MappingInstruction
}

func (m *AlmanacMapper) findNext(item int) int {
	for _, instruction := range m.mappings {
		slog.Debug("Instructions to check", "instruction", instruction)
		mappedItem := instruction.findNext(item)
		if mappedItem != item {
			return mappedItem
		}
	}
	return item
}

type MappingInstruction struct {
	destStart   int
	sourceStart int
	rangeLength int
}

func (i *MappingInstruction) findNext(item int) int {
	var diff int

	if item >= i.sourceStart && item < i.sourceStart+i.rangeLength {
		diff = i.destStart - i.sourceStart
		slog.Debug("Item Found in Mapping", "item", item, "destination", item+diff, "mapping", i)
	} else {
		diff = 0
	}

	return item+diff
}

// Parsing Logic
type AlmanacParser struct {
	scanner bufio.Scanner
}

// func createAlmanacParser(file *os.File, newline byte) *AlmanacParser {
// 	p := &AlmanacParser{}
// 	p.scanner = *bufio.NewScanner(file)
// 	split := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
// 		if atEOF && len(data) == 0 {
// 			return 0, nil, nil
// 		}

// 		if atEOF {
// 			return len(data), data, nil
// 		}

// 		for i := 0; i < (cap(data) - 1); i++ {
// 			if bytes.Equal(data[i:i+2], []byte{newline, newline}) {
// 				return i+2, data[:i+2], nil
// 			}
// 		}

// 		return 0, nil, nil
// 	}
// 	p.scanner.Split(split)
// 	return p
// }

func createAlmanacParser(file *os.File) *AlmanacParser {
	p := &AlmanacParser{}
	p.scanner = *bufio.NewScanner(file)
	return p
}

func (p *AlmanacParser) createAlamanac() *Almanac {
	a := &Almanac{}
	a.seeds = p.getSeeds()
	a.seedToSoil = p.getNextMapping()
	a.soilToFertilizer = p.getNextMapping()
	a.fertilizerToWater = p.getNextMapping()
	a.waterToLight = p.getNextMapping()
	a.lightToTemperature = p.getNextMapping()
	a.temperatureToHumidity = p.getNextMapping()
	a.humidityToLocation = p.getNextMapping()
	return a
}

func (p *AlmanacParser) getNextTokens() []string {
	tokens := []string{}

	for p.scanner.Scan() {
		line := p.scanner.Text()
		if len(line) == 0 {
			break
		}
		line = strings.TrimSuffix(line, "\n")
		tokens = append(tokens, line)
	}

	return tokens
}

func (p *AlmanacParser) getSeeds() []int {
	seedLine := p.getNextTokens()
	if len(seedLine) != 1 {
		log.Fatalf("Invalid Header %v\n", seedLine)
	}
	seedLine = strings.Split(seedLine[0], ":")
	strSeeds := strings.Fields(seedLine[1])
	seeds := utils.StringSliceToIntSlice(strSeeds)
	return seeds
}

func (p *AlmanacParser) getNextMapping() AlmanacMapper {
	strMapping := p.getNextTokens()
	cMappings := []MappingInstruction{}
	for _, aMapping := range strMapping[1:] {
		aFields := strings.Fields(aMapping)
		ints := utils.StringSliceToIntSlice(aFields)
		c := MappingInstruction{ints[0], ints[1], ints[2]}
		cMappings = append(cMappings, c)
	}
	return AlmanacMapper{cMappings}
}
