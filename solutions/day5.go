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

	location := almanac.FindLocationFromRange()
	fmt.Printf("Next here: %v\n", location)
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

func (a *Almanac) FindLocationFromRange() int {
	minimum := 1024*1024*1024*1024

	locations := a.SeedRanges()
	for _, mapper := range a.mappers() {
		nextLocations := []SeedRange{}
		for _, location := range locations {
			result := mapper.findNextRanges(location)
			nextLocations = append(nextLocations, result...)
		}
		locations = nextLocations
	}

	minimums := []int{}
	for _, seed := range locations {
		minimums = append(minimums, seed.Start)
	}

	for _, min := range minimums {
		if min < minimum {
			minimum = min
		}
	}

	return minimum
}

func (a *Almanac) SeedRanges() []SeedRange {
	ranges := []SeedRange{}

	for i := 0; i < len(a.seeds); i = i + 2 {
		ranges = append(ranges, SeedRange{a.seeds[i], a.seeds[i+1]})
	}

	return ranges
}

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

func (m *AlmanacMapper) findNextRanges(seeds SeedRange) []SeedRange {
	var ranges []SeedRange

	for _, instruction := range m.mappings {
		ranges = instruction.transform(seeds)
		if !(len(ranges) == 1 && ranges[0] == seeds) {
			slog.Debug("Transformed range", "range", ranges, "before", seeds, "instruction", instruction)
			return ranges
		}
	}

	return ranges
}

type MappingInstruction struct {
	Destination int
	Source      int
	Length      int
}

func (i *MappingInstruction) Transform() int {
	return i.Destination - i.Source
}

func (i *MappingInstruction) findNext(item int) int {
	var diff int

	if item >= i.Source && item < i.Source+i.Length {
		diff = i.Destination - i.Source
		slog.Debug("Item Found in Mapping", "item", item, "destination", item+diff, "mapping", i)
	} else {
		diff = 0
	}

	return item + diff
}

func (i *MappingInstruction) transform(r SeedRange) []SeedRange {
	if r.Start >= i.Source && r.Start+r.Length <= i.Source+i.Length { // Full Match
		dest := SeedRange{Start: r.Start + i.Transform(), Length: r.Length}
		return []SeedRange{dest}
	} else if r.Start < i.Source && r.Start+r.Length > i.Source+i.Length { // 3-Part
		pre := SeedRange{r.Start, i.Source - r.Start - 1}
		dest := SeedRange{r.Start + i.Transform(), i.Length}
		post := SeedRange{i.Source + i.Length + 1, (r.Start+r.Length) - (i.Source+i.Length) - 1}
		return []SeedRange{pre, dest, post}
	} else if (r.Start >= i.Source && r.Start+r.Length <= i.Source+i.Length) && r.Start+r.Length > i.Source+i.Length { // R-Part
		dest := SeedRange{r.Start + i.Transform(), i.Source + i.Length - r.Start}
		post := SeedRange{i.Source + i.Length + 1, (r.Start + r.Length) - (i.Source + i.Length) - 1}
		return []SeedRange{dest, post}
	} else if r.Start < i.Source && (r.Start+r.Length <= i.Source+i.Length && r.Start+r.Length >= i.Source) { // L-Part
		pre := SeedRange{r.Start, i.Source - r.Start - 1}
		dest := SeedRange{i.Source + i.Transform(), (r.Start + r.Length) - i.Source}
		return []SeedRange{pre, dest}
	} else { // No Match
		return []SeedRange{r}
	}
}

type SeedRange struct {
	Start  int
	Length int
}

// Parsing Logic
type AlmanacParser struct {
	scanner bufio.Scanner
}

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
