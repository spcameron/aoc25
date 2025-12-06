package main

import (
	"bufio"
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

	operandLines := [][]string{}
	operators := []string{}

	for _, line := range lines {
		if len(line) == 0 {
			continue
		}

		parsedLine, err := parseLine(line)
		if err != nil {
			panic(err)
		}

		if parsedLine[0] == "+" || parsedLine[0] == "*" {
			operators = parsedLine
			continue
		}

		operandLines = append(operandLines, parsedLine)
	}

	if len(operandLines[0]) != len(operators) {
		panic("the operands and operators are different lengths")
	}

	total := 0
	for i, operator := range operators {
		operands := []string{}
		for _, line := range operandLines {
			operands = append(operands, line[i])
		}

		switch operator {
		case "+":
			result, err := add(operands)
			if err != nil {
				panic(err)
			}

			total += result
		case "*":
			result, err := multiply(operands)
			if err != nil {
				panic(err)
			}

			total += result
		default:
			panic("unrecognized operator")
		}
	}

	fmt.Printf("the sum of each operation result is: %d\n", total)

}

func parseLine(input []byte) ([]string, error) {
	scanner := bufio.NewScanner(bytes.NewReader(input))
	scanner.Split(bufio.ScanWords)

	output := []string{}
	for scanner.Scan() {
		output = append(output, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return output, nil
}

func add(operands []string) (int, error) {
	sum := 0
	for _, operand := range operands {
		i, err := strconv.Atoi(operand)
		if err != nil {
			return 0, err
		}

		sum += i
	}

	return sum, nil
}

func multiply(operands []string) (int, error) {
	product := 1
	for _, operand := range operands {
		i, err := strconv.Atoi(operand)
		if err != nil {
			return 0, err
		}

		product *= i
	}

	return product, nil
}
