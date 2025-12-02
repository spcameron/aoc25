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

	for m := 1; m <= n/2; m++ {
		if isRepeatedPattern(id, m) {
			return true
		}
	}

	return false
}

func isRepeatedPattern(s string, patternLength int) bool {
	if patternLength <= 0 || len(s)%patternLength != 0 {
		return false
	}

	pattern := s[:patternLength]

	for pos := patternLength; pos < len(s); pos += patternLength {
		if s[pos:pos+patternLength] != pattern {
			return false
		}
	}

	return true
}

// func invalidID(id string) bool {
// 	n := len(id)
// 	for m := 1; m <= n/2; m++ {
// 		if n%m != 0 {
// 			continue
// 		}
//
// 		p := id[0:m]
// 		for j := m; j <= n-m; j += m {
// 			if p != id[j:j+m] {
// 				break
// 			}
//
// 			if j+m == n {
// 				return true
// 			}
// 		}
//
// 	}
//
// 	return false
// }

func sumInvalidIDs(invalidIDs []int) int {
	total := 0
	for _, id := range invalidIDs {
		total += id
	}

	return total
}
