package main

import (
	"bytes"
	"fmt"
	"os"
	"slices"
	"strconv"
	"time"
)

func main() {
	start := time.Now()

	path := "./input.txt"
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	lines := bytes.Split(data, []byte("\n"))
	machines := []Machine{}
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}

		machine, err := parseLine(line)
		if err != nil {
			panic(err)
		}

		machines = append(machines, *machine)
	}

	totalPresses := 0
	for i, machine := range machines {
		fmt.Printf("processing machine #%d\n", i+1)
		presses := machine.Configure()
		totalPresses += presses
	}

	elapsed := time.Since(start).Seconds()
	fmt.Printf("total presses to configure all machines: %d\n", totalPresses)
	fmt.Printf("time elapsed: %vs\n", elapsed)
}

type Machine struct {
	Indicator Indicator
	Buttons   []Button
	Joltage   Joltage
}

func (m Machine) Configure() int {
	n := 1

	for {
		q := Queue{
			cache: make(map[[2]string]bool),
		}

		startNode := Node{make([]int, len(m.Joltage.Indexes)), n}

		q.Push(startNode)

		for !q.isEmpty() {
			node := q.Pop()
			skip := false

			for i := range len(m.Joltage.Indexes) {
				if node.joltageVals[i] > m.Joltage.Indexes[i] {
					skip = true
				}
			}

			if skip {
				continue
			}

			if m.AttemptConfigure(node.joltageVals, node.depth, &q) {
				return node.depth
			}
		}

		n++
	}
}

func (m Machine) AttemptConfigure(curr []int, n int, q *Queue) bool {
	fmt.Printf("attempting configuration with curr: %v, n: %d\n", curr, n)
	if n == 0 {
		return false
	}

	for _, button := range m.Buttons {
		val := button.Press(curr)
		if slices.Equal(val, m.Joltage.Indexes) {
			return true
		}

		valString := fmt.Sprint(val)
		depthString := fmt.Sprint(n - 1)
		if !q.cache[[2]string{valString, depthString}] {
			q.Push(Node{val, n - 1})
			q.cache[[2]string{valString, depthString}] = true
		}
	}

	return false
}

type Indicator struct {
	RawInput string
}

type Button struct {
	Indexes []int
}

func (b Button) Press(input []int) []int {
	output := make([]int, len(input))
	copy(output, input)
	for _, idx := range b.Indexes {
		output[idx]++
	}

	return output
}

type Joltage struct {
	Indexes []int
}

func parseLine(input []byte) (*Machine, error) {
	tokens := bytes.Split(input, []byte(" "))

	n := len(tokens)
	indicatorToken := tokens[0]
	buttonTokens := tokens[1 : n-1]
	joltageToken := tokens[n-1]

	indicator, err := parseIndicatorToken(indicatorToken)
	if err != nil {
		return nil, err
	}

	buttons := []Button{}
	for _, buttonToken := range buttonTokens {
		button, err := parseButtonToken(buttonToken)
		if err != nil {
			return nil, err
		}

		buttons = append(buttons, button)
	}

	joltage, err := parseJoltageToken(joltageToken)
	if err != nil {
		return nil, err
	}

	return &Machine{indicator, buttons, joltage}, nil
}

func parseIndicatorToken(token []byte) (Indicator, error) {
	n := len(token)
	if token[0] != '[' || token[n-1] != ']' {
		return Indicator{}, fmt.Errorf("invalid indicator token: %s", string(token))
	}

	s := string(token)

	return Indicator{s}, nil
}

func parseButtonToken(token []byte) (Button, error) {
	n := len(token)
	if token[0] != '(' || token[n-1] != ')' {
		return Button{}, fmt.Errorf("invalid button token: %s", string(token))
	}

	indexes := []int{}
	for _, c := range token {
		if c >= '0' && c <= '9' {
			n, err := strconv.Atoi(string(c))
			if err != nil {
				return Button{}, err
			}

			indexes = append(indexes, n)
		}
	}

	return Button{indexes}, nil
}

func parseJoltageToken(token []byte) (Joltage, error) {
	n := len(token)
	if token[0] != '{' || token[n-1] != '}' {
		return Joltage{}, fmt.Errorf("invalid joltage token: %s", string(token))
	}

	token = token[1 : n-1]
	valBytes := bytes.Split(token, []byte(","))

	values := []int{}
	for _, b := range valBytes {
		n, err := strconv.Atoi(string(b))
		if err != nil {
			return Joltage{}, err
		}

		values = append(values, n)
	}

	return Joltage{values}, nil
}

type Queue struct {
	slice []Node
	head  int
	cache map[[2]string]bool
}

func (q *Queue) Push(n Node) {
	q.slice = append(q.slice, n)
}

func (q *Queue) Pop() Node {
	n := q.slice[q.head]
	q.head++
	return n
}

func (q *Queue) isEmpty() bool {
	return q.head == len(q.slice)
}

type Node struct {
	joltageVals []int
	depth       int
}
