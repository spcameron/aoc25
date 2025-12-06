package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
)

func main() {
	path := "./input.txt"
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	lines := bytes.Split(data, []byte("\n"))
	for len(lines[len(lines)-1]) == 0 {
		lines = lines[:len(lines)-1]
	}

	total := 0
	columns := len(lines[0])

	digits := [][]byte{}
	operands := []int{}
	operator := ""
	spaces := 0
	for c := range columns {
		spaces = 0
		digits = digits[:0]
		for _, line := range lines {
			switch line[c] {
			case ' ':
				spaces += 1
			case '+':
				operator = "+"
			case '*':
				operator = "*"
			default:
				digits = append(digits, []byte{line[c]})
			}
		}

		if spaces == len(lines) {
			total += operate(operands, operator)

			operands = operands[:0]
			operator = ""
			continue
		}

		operand, err := processDigits(digits)
		if err != nil {
			panic(err)
		}

		operands = append(operands, operand)

		if c == columns-1 {
			total += operate(operands, operator)
		}
	}

	fmt.Printf("the sum of each operation result is: %d\n", total)
}

func processDigits(digits [][]byte) (int, error) {
	stringVal := string(bytes.Join(digits, []byte("")))
	i, err := strconv.Atoi(stringVal)
	if err != nil {
		return 0, err
	}

	return i, nil
}

func operate(operands []int, operator string) int {
	result := 0
	if operator == "*" {
		result = 1
	}

	for _, operand := range operands {
		switch operator {
		case "+":
			result += operand
		case "*":
			result *= operand
		}
	}

	return result
}
