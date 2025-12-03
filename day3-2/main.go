package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
)

func main() {
	path := "./day3-puzzle-input.txt"
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	lines := bytes.Split(data, []byte("\n"))

	joltage, err := maxJoltage(lines, 12)
	if err != nil {
		panic(err)
	}

	fmt.Printf("the total maximum joltage: %d\n", joltage)
}

func maxJoltage(lines [][]byte, batteries int) (int, error) {
	total := 0

	for _, line := range lines {
		i, err := parseLine(line, batteries)
		if err != nil {
			return 0, err
		}

		total += i
	}

	return total, nil
}

func parseLine(input []byte, n int) (int, error) {
	if len(input) == 0 {
		return 0, nil
	}

	digits := make([]byte, 0, n)

	leftBoundaryIndex := 0
	for i := range n {
		reserve := n - i - 1
		rightBoundaryIndex := len(input) - reserve
		val, idx := findLeftmostHighNumber(input[leftBoundaryIndex:rightBoundaryIndex])
		digits = append(digits, val)
		leftBoundaryIndex += idx + 1
	}

	maxJoltage, err := strconv.Atoi(string(digits))
	if err != nil {
		return 0, err
	}
	return maxJoltage, nil
}

func findLeftmostHighNumber(input []byte) (val byte, idx int) {
	val = '0'
	idx = -1

	for i, v := range input {
		if v == '9' {
			val = v
			idx = i
			break
		}

		if v > val {
			val = v
			idx = i
		}
	}

	return val, idx
}
