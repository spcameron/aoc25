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

	total := 0

	for _, line := range lines {
		i, err := parseLine(line)
		if err != nil {
			panic(err)
		}
		total += i
	}

	fmt.Printf("the total maximum joltage: %d\n", total)
}

func parseLine(input []byte) (int, error) {
	if len(input) == 0 {
		return 0, nil
	}

	firstNumber, firstIdx := findLeftmostHighNumber(input[:len(input)-1])
	secondNumber, _ := findLeftmostHighNumber(input[firstIdx+1:])
	combinedNumber := []byte{firstNumber, secondNumber}

	fmt.Printf("first: %v, second: %v\n", firstNumber, secondNumber)

	maxJoltage, err := strconv.Atoi(string(combinedNumber))
	if err != nil {
		return 0, err
	}

	fmt.Printf("combined: %d\n", maxJoltage)

	return maxJoltage, nil
}

func findLeftmostHighNumber(input []byte) (val byte, idx int) {
	val = '0'
	idx = -1

	for i, v := range input {
		if v > val {
			val = v
			idx = i
		}
	}

	return val, idx
}
