package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/ryanpdenoux/advent-of-code/utils"
	"github.com/ryanpdenoux/advent-of-code/solutions"
)

// register solution funcs
var solutionMap = map[int]func(*os.File){
	1: solutions.Day1,
	2: solutions.Day2,
	3: solutions.Day3,
	4: solutions.Day4,
	5: solutions.Day5,
	6: solutions.Day6,
	7: solutions.Day7,
	8: solutions.Day8,
	9: solutions.Day9,
}

var (
	day = flag.Int("day",
		0,
		"day[n] to select for solution",
	)
	year = flag.Int("year",
		2023,
		"Year that solution exists",
	)
	debug = flag.Bool("debug",
		false,
		"Toggle Debugging logs for solution",
	)
	quiet = flag.Bool("quiet",
		false,
		"Turn off logging output",
	)
)

func pickDay() (int) {
	var day int

	fmt.Print("Pick a day: ")
	fmt.Scan(&day)

	return day
}

func setupLogging(debug, quiet bool) {
	opts := &slog.HandlerOptions{}
	if debug {
		opts.Level = slog.LevelDebug
	}
	if quiet {
		opts.Level = slog.LevelError
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, opts))
	slog.SetDefault(logger)
}

func main() {
	flag.Parse()

	setupLogging(*debug, *quiet)

	if *day == 0 {
		*day = pickDay()
	}

	solution_func := solutionMap[*day]
	dataPath := utils.BuildDataPath(*day)
	file, err := os.Open(dataPath)
	if err != nil {
		log.Fatalf("Could not open file %v: %v", dataPath, err)
	}
	defer file.Close()

	solution_func(file)
}
