package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
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
	vertices := getVertices(lines)

	maxArea := 0
	for i := 0; i < len(vertices); i++ {
		for j := i + 1; j < len(vertices); j++ {
			area := Rectangle{vertices[i], vertices[j]}.Area()
			if maxArea < area {
				maxArea = area
			}
		}
	}

	elapsed := time.Since(start).Seconds()
	fmt.Printf("the largest area possible is: %d\n", maxArea)
	fmt.Printf("time elapsed: %vs\n", elapsed)
}

type Point struct {
	x int
	y int
}

type Rectangle struct {
	v1 Point
	v2 Point
}

func (r Rectangle) Area() int {
	dx := r.v1.x - r.v2.x
	if dx < 0 {
		dx = -dx
	}

	dy := r.v1.y - r.v2.y
	if dy < 0 {
		dy = -dy
	}

	return dx * dy
}

func getVertices(lines [][]byte) []Point {
	points := []Point{}
	for _, line := range lines {
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

		points = append(points, Point{x, y})
	}

	return points
}
