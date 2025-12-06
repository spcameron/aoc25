package main

import (
	"bytes"
	"fmt"
	"os"
)

func main() {
	path := "../input.txt"
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

	operands := []int{}
	operator := ""
	for c := range columns {
		spaces := 0
		value := 0
		hasDigit := false
		for _, line := range lines {
			switch line[c] {
			case ' ':
				spaces += 1
			case '+':
				operator = "+"
			case '*':
				operator = "*"
			default:
				ch := line[c]
				if ch >= '0' && ch <= '9' {
					value = value*10 + int(ch-'0')
					hasDigit = true
				}
			}
		}

		if hasDigit {
			operands = append(operands, value)
		}

		if spaces == len(lines) || c == columns-1 {
			if len(operands) == 0 || operator == "" {
				panic("malformed puzzle input")
			}

			total += operate(operands, operator)

			operands = operands[:0]
			operator = ""
			continue
		}

	}

	fmt.Printf("the sum of each operation result is: %d\n", total)
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
