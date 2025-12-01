package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
)

func main() {
	path := "./day1-puzzle-input.txt"
	data, err := readFile(path)
	if err != nil {
		panic(err)
	}

	lines := splitLines(data)
	start := 50

	pass, err := countZeroHits(start, lines)
	if err != nil {
		panic(err)
	}
	fmt.Printf("password: %d", pass)
}

func readFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func splitLines(data []byte) [][]byte {
	newlineSeparator := []byte("\n")
	return bytes.Split(data, newlineSeparator)
}

func parseTurn(line []byte) (int, string, error) {
	direction := string(line[:1])
	if direction != "R" && direction != "L" {
		return 0, "", fmt.Errorf("unrecognized turn direction %s", direction)
	}

	distance := line[1:]
	s := string(distance)
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, "", err
	}

	return i, direction, nil
}

func countZeroHits(start int, turns [][]byte) (int, error) {
	pos := start
	count := 0
	for _, turn := range turns {
		if len(turn) == 0 {
			continue
		}

		i, dir, err := parseTurn(turn)
		if err != nil {
			return 0, err
		}

		hits := 0
		pos, hits = updatePosition(pos, i, dir)
		count += hits
	}

	return count, nil
}

func updatePosition(curr, distance int, direction string) (int, int) {
	shift := 1
	if direction == "L" {
		shift = -1
	}

	hits := 0

	for range distance {
		curr += shift
		if curr == -1 {
			curr = 99
		}
		if curr == 100 {
			curr = 0
		}
		if curr == 0 {
			hits += 1
		}
	}

	return curr, hits
}
