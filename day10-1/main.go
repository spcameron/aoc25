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
}

func (m Machine) Start() int {
	n := 1
	for {
		if m.AttemptStart(0, n) {
			break
		}
		n++
	}

	return n
}

func (m Machine) AttemptStart(curr int64, n int) bool {
	if n == 0 {
		return false
	}

	for _, button := range m.Buttons {
		val := button.Press(curr)
		if val == m.Indicator.BinaryValue {
			return true
		}

		if m.AttemptStart(val, n-1) {
			return true
		}
	}

	return false
}

type Indicator struct {
	BinaryString string
	BinaryValue  int64
}

func (i Indicator) Length() int {
	return len(i.BinaryString)
}

type Button struct {
	BinaryString string
	BinaryValue  int64
}

func (b Button) Press(input int64) int64 {
	return input ^ b.BinaryValue
}

type Joltage struct {
	RawInput string
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
		button, err := parseButtonToken(buttonToken, indicator.Length())
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
	binaryVal, err := strconv.ParseInt(binaryString, 2, 0)
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

	places := []int{}
	for _, c := range token {
		if c >= '0' && c <= '9' {
			n, err := strconv.Atoi(string(c))
			if err != nil {
				return Button{}, err
			}

			places = append(places, n)
		}
	}

	var b strings.Builder
	idx := 0
	for i := 0; i < indicatorLength; i++ {
		if idx >= len(places) || i != places[idx] {
			b.WriteString("0")
			continue
		}

		b.WriteString("1")
		idx++
	}

	binaryString := b.String()
	binaryVal, err := strconv.ParseInt(binaryString, 2, 0)
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
