package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	start := time.Now()

	path := "../input.txt"
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

		machines = append(machines, machine)
	}

	totalPresses := 0
	for i, machine := range machines {
		fmt.Printf("processing machine #%d\n", i+1)
		presses := machine.Start()
		totalPresses += presses
	}

	elapsed := time.Since(start).Seconds()
	fmt.Printf("total presses to start all machines: %d\n", totalPresses)
	fmt.Printf("time elapsed: %vs\n", elapsed)
}

type Machine struct {
	Indicator Indicator
	Buttons   []Button
	Joltage   Joltage
	Cache     map[CacheKey]int
}

type CacheKey struct {
	i    int
	mask uint64
}

func (m *Machine) Start() int {
	return m.dp(0, 0)
}

func (m *Machine) dp(i int, mask uint64) int {
	if v, ok := m.Cache[CacheKey{i, mask}]; ok {
		return v
	}

	if mask == m.Indicator.BinaryValue {
		m.Cache[CacheKey{i, mask}] = 0
		return 0
	}

	if i == len(m.Buttons) {
		return len(m.Buttons) + 1
	}

	skip := m.dp(i+1, mask)
	use := 1 + m.dp(i+1, m.Buttons[i].Press(mask))

	result := min(skip, use)
	m.Cache[CacheKey{i, mask}] = result
	return result
}

type Indicator struct {
	BinaryString string
	BinaryValue  uint64
}

func (i Indicator) Length() int {
	return len(i.BinaryString)
}

type Button struct {
	BinaryString string
	BinaryValue  uint64
}

func (b Button) Press(input uint64) uint64 {
	return input ^ b.BinaryValue
}

type Joltage struct {
	RawInput string
}

func parseLine(input []byte) (Machine, error) {
	tokens := bytes.Split(input, []byte(" "))

	n := len(tokens)
	indicatorToken := tokens[0]
	buttonTokens := tokens[1 : n-1]
	joltageToken := tokens[n-1]

	indicator, err := parseIndicatorToken(indicatorToken)
	if err != nil {
		return Machine{}, err
	}

	buttons := []Button{}
	for _, buttonToken := range buttonTokens {
		button, err := parseButtonToken(buttonToken, indicator.Length())
		if err != nil {
			return Machine{}, err
		}

		buttons = append(buttons, button)
	}

	joltage, err := parseJoltageToken(joltageToken)
	if err != nil {
		return Machine{}, err
	}

	cache := make(map[CacheKey]int)

	return Machine{indicator, buttons, joltage, cache}, nil
}

func parseIndicatorToken(token []byte) (Indicator, error) {
	n := len(token)
	if token[0] != '[' || token[n-1] != ']' {
		return Indicator{}, fmt.Errorf("invalid indicator token: %s", string(token))
	}

	var b strings.Builder
	for _, c := range token {
		switch c {
		case '.':
			b.WriteString("0")
		case '#':
			b.WriteString("1")
		default:
			continue
		}
	}

	binaryString := b.String()
	binaryVal, err := strconv.ParseUint(binaryString, 2, 0)
	if err != nil {
		return Indicator{}, err
	}

	return Indicator{binaryString, binaryVal}, nil
}

func parseButtonToken(token []byte, indicatorLength int) (Button, error) {
	n := len(token)
	if token[0] != '(' || token[n-1] != ')' {
		return Button{}, fmt.Errorf("invalid button token: %s", string(token))
	}

	indexes := []int{}
	for _, c := range token {
		if c >= '0' && c <= '9' {
			v, err := strconv.Atoi(string(c))
			if err != nil {
				return Button{}, err
			}

			indexes = append(indexes, v)
		}
	}

	var b strings.Builder
	idx := 0
	for i := 0; i < indicatorLength; i++ {
		if idx >= len(indexes) || i != indexes[idx] {
			b.WriteString("0")
			continue
		}

		b.WriteString("1")
		idx++
	}

	binaryString := b.String()
	binaryVal, err := strconv.ParseUint(binaryString, 2, 0)
	if err != nil {
		return Button{}, err
	}

	return Button{binaryString, binaryVal}, nil
}

func parseJoltageToken(token []byte) (Joltage, error) {
	n := len(token)
	if token[0] != '{' || token[n-1] != '}' {
		return Joltage{}, fmt.Errorf("invalid joltage token: %s", string(token))
	}

	s := string(token)

	return Joltage{s}, nil
}
