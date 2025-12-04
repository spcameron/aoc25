package main

import (
	"bytes"
	"fmt"
	"os"
)

func main() {
	path := "./day4-puzzle-input.txt"
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	data = bytes.TrimSpace(data)

	matrix := bytes.Split(data, []byte("\n"))
	matrixHeight := len(matrix)
	matrixWidth := len(matrix[0])

	totalRolls := 0

	for y := range matrixHeight {
		for x := range matrixWidth {
			if !isPaperRoll(matrix[y][x]) {
				continue
			}

			adjacentRolls := countAdjacentRolls(x, y, matrix)

			if adjacentRolls < 4 {
				totalRolls++
			}
		}
	}

	fmt.Printf("total rolls that can be moved: %d\n", totalRolls)
}

func isPaperRoll(input byte) bool {
	return input == '@'
}

func isValidPositionAndPaperRoll(x, y int, matrix [][]byte) bool {
	if x < 0 || y < 0 || x >= len(matrix[0]) || y >= len(matrix) {
		return false
	}

	if !isPaperRoll(matrix[y][x]) {
		return false
	}

	return true
}

func countAdjacentRolls(x, y int, matrix [][]byte) int {
	adjacentRolls := 0

	for i := x - 1; i <= x+1; i++ {
		for j := y - 1; j <= y+1; j++ {
			if i == x && j == y {
				continue
			}

			if isValidPositionAndPaperRoll(i, j, matrix) {
				adjacentRolls++
			}
		}
	}

	return adjacentRolls
}
