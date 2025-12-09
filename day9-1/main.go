package main

import (
	"bytes"
	"fmt"
	"math"
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
	redTiles := gedRedTiles(lines)

	maxArea := 0
	for i := 0; i < len(redTiles); i++ {
		for j := i + 1; j < len(redTiles); j++ {
			currArea := redTiles[i].area(redTiles[j])
			if maxArea < currArea {
				maxArea = currArea
			}
		}
	}

	fmt.Printf("the largest area possible is: %d\n", maxArea)
}

type point struct {
	x, y int
}

func (p point) area(q point) int {
	dx := math.Abs(float64(p.x-q.x)) + 1
	dy := math.Abs(float64(p.y-q.y)) + 1

	return int(dx * dy)
}

func gedRedTiles(input [][]byte) []point {
	points := []point{}
	for _, line := range input {
		if len(line) == 0 {
			continue
		}

		coords := bytes.Split(line, []byte(","))
		if len(coords) != 2 {
			panic(fmt.Sprintf("unexpected line format: %s", string(line)))
		}

		x, err := strconv.Atoi(string(coords[0]))
		if err != nil {
			panic(err)
		}

		y, err := strconv.Atoi(string(coords[1]))
		if err != nil {
			panic(err)
		}

		points = append(points, point{x, y})
	}

	return points
}
