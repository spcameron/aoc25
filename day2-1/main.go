package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	path := "./day2-puzzle-input.txt"

	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	invalidIDs := make([]int, 0)

	for part := range bytes.SplitSeq(data, []byte(",")) {
		trimmedPart := strings.TrimSpace(string(part))
		start, end, err := getIDRangeBounds(trimmedPart)
		if err != nil {
			panic(err)
		}

		invalidIDs = append(invalidIDs, validateIDs(start, end)...)
	}

	sum := sumInvalidIDs(invalidIDs)

	fmt.Printf("total of the invalid IDs: %d\n", sum)
}

func getIDRangeBounds(input string) (start, end int, err error) {
	parts := strings.Split(input, "-")
	if len(parts) != 2 {
		return 0, 0, errors.New("invalid product ID range encountered")
	}

	start, err = strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, err
	}

	end, err = strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, err
	}

	return start, end, nil
}

func validateIDs(start, end int) []int {
	invalidIDs := make([]int, 0)
	for i := start; i <= end; i++ {
		if invalidID(strconv.Itoa(i)) {
			invalidIDs = append(invalidIDs, i)
		}
	}

	return invalidIDs
}

func invalidID(id string) bool {
	n := len(id)
	if n%2 != 0 {
		return false
	}

	firstHalf := id[:n/2]
	secondHalf := id[n/2:]

	if firstHalf != secondHalf {
		return false
	}

	return true
}

func sumInvalidIDs(invalidIDs []int) int {
	total := 0
	for _, id := range invalidIDs {
		total += id
	}

	return total
}
