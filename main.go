package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/ryanpdenoux/advent-of-code/utils"
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
)

func pickDay() (int) {
	var day int

	fmt.Print("Pick a day: ")
	fmt.Scan(&day)

	return day
}

func setupLogging(debug bool) {
	opts := &slog.HandlerOptions{}
	if debug {
		opts.Level = slog.LevelDebug
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, opts))
	slog.SetDefault(logger)
}

func main() {
	flag.Parse()

	setupLogging(*debug)

	if *day == 0 {
		*day = pickDay()
	}

	solution_func := solutionMap[*day]
	dataPath := utils.BuildDataPath(*year, *day)
	file, err := os.Open(dataPath)
	if err != nil {
		log.Fatalf("Could not open file %v: %v", dataPath, err)
	}
	defer file.Close()

	solution_func(file)
}
