package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ryanpdenoux/advent-of-code/2023"
)

// register solution funcs
var solutionMap = map[int]func(*os.File){
	1: solution.Day1,
	2: solution.Day2,
	3: solution.Day3,
}

var (
	day = flag.Int("day",
		0,
		"day[n] to select for solution")
	year = flag.Int("year",
		2023,
		"Year that solution exists")
)

func pickDay() (int) {
	var day int

	fmt.Print("Pick a day: ")
	fmt.Scan(&day)

	return day
}

func buildDataPath(day int) (string) {
	dataPath := fmt.Sprintf("%d/day%d-input.txt", *year, day)
	return dataPath
}

func main() {
	flag.Parse()

	if *day == 0 {
		*day = pickDay()
	}

	solution_func := solutionMap[*day]
	dataPath := buildDataPath(*day)
	file, err := os.Open(dataPath)
	if err != nil {
		log.Fatalf("Could not open file %v: %v", dataPath, err)
	}
	solution_func(file)
	file.Close()
}
