package solution

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"unicode/utf8"
)

func Day1(file *os.File) {
	sum, err := sumCalibrationValues(file)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Sum of values: %d\n", sum)
}

var parsingFuncs = []func([]rune, int) (int, bool){
	strToInt,
}

func strToInt(chars []rune, i int) (int, bool) {
	num, err := strconv.Atoi(string(chars[i]))
	if err != nil {
		return 0, false
	}

	return num, true
}

// func completeTrie(chars []rune, i int) (int, bool) {}

func sumCalibrationValues(file *os.File) (int, error) {
	scanner := bufio.NewScanner(file)
	var sum int

	for scanner.Scan() {
		line := scanner.Text()
		value, err := constructValue(line)
		if err != nil {
			return 0, fmt.Errorf("Something busted: %v", err)
		}
		sum += value
	}

	return sum, nil
}

func constructValue(line string) (int, error) {
	first, ok := pickFirstDigit(line)
	if ok != true {
		return 0, fmt.Errorf("No digits in string: %v", line)
	}

	last, _ := pickFirstDigit(reverseString(line))
	value := ((10 * first) + last)
	return value, nil
}

func reverseString(input string) string {
	r := make([]byte, 0, len(input))
	for len(input) > 0 {
		_, n := utf8.DecodeLastRuneInString(input)
		i := len(input) - n
		r = append(r, input[i:]...)
		input = input[:i]
	}
	return string(r)
}

func pickFirstDigit(input string) (int, bool) {
	chars := []rune(input)

	for pos := range(chars) {
		for _, parsingFn := range(parsingFuncs) {
			num, ok := parsingFn(chars, pos)
			if ok == true {
				return num, ok
			}
		}
	}

	return 0, false
}
