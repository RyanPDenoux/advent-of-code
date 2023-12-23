package utils

import (
	"fmt"
	"math"
	"strconv"
)

func BuildDataPath(day int) (string) {
	dataPath := fmt.Sprintf("solutions/day%d-input.txt", day)
	return dataPath
}

func IsLetter(char byte) bool {
	return 'a' <= char && char <= 'z' || 'A' <= char && char <= 'Z'
}

func IsDigit(char byte) bool {
	return '0' <= char && char <= '9'
}

func MakeRange(min, max int) []int {
	a := make([]int, max-min)
	for i := range a {
		a[i] = min + i
	}
	return a
}

func LengthOfInt(number int) int {
	return int(math.Log10(float64(number))) + 1
}

func FindNumberInBytes(data []byte, pos int) (int, bool) {
	var num []byte

	// check length to avoid index errors
	for i := pos; i < len(data) && IsDigit(data[i]); i++ {
		num = append(num, data[i])
	}

	intNum, err := strconv.Atoi(string(num))
	if err != nil {
		return 0, false
	}

	return intNum, true
}
