package utils

import (
	"log"
	"strconv"
)

func StringSliceToIntSlice(s []string) []int {
	ints := []int{}

	for _, char := range s {
		i, err := strconv.Atoi(char)
		if err != nil {
			log.Fatalf("Not all elements can be converted to int: %v %v", s, err)
		}
		ints = append(ints, i)
	}

	return ints
}
