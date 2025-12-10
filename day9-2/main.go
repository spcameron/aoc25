package main

import (
	"bytes"
	"fmt"
	"math"
	"os"
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
	redTiles := getRedTiles(lines)

	vEdges, hEdges := constructEdges(redTiles)

	maxArea := 0
	for i := 0; i < len(redTiles); i++ {
		for j := i + 1; j < len(redTiles); j++ {

			if !validCandidateRectangle(redTiles[i], redTiles[j], redTiles, vEdges, hEdges) {
				continue
			}

			currArea := redTiles[i].area(redTiles[j])
			if maxArea < currArea {
				maxArea = currArea
			}
		}
	}

	elapsed := time.Since(start).Seconds()
	fmt.Printf("the largest area possible is: %d\n", maxArea)
	fmt.Printf("time elapsed: %vs\n", elapsed)
}

type point struct {
	x, y int
}

func (p point) area(q point) int {
	dx := math.Abs(float64(p.x-q.x)) + 1
	dy := math.Abs(float64(p.y-q.y)) + 1

	return int(dx * dy)
}

type vEdge struct {
	x, y1, y2 int
}

type hEdge struct {
	y, x1, x2 int
}

func constructEdges(vertices []point) (vEdges []vEdge, hEdges []hEdge) {
	for i := range len(vertices) {
		a := vertices[i]
		b := vertices[(i+1)%len(vertices)]

		if a.x == b.x {
			y1 := a.y
			y2 := b.y
			if y1 > y2 {
				y1, y2 = y2, y1
			}

			vEdges = append(vEdges, vEdge{
				x:  a.x,
				y1: y1,
				y2: y2,
			})
		}

		if a.y == b.y {
			x1 := a.x
			x2 := b.x
			if x1 > x2 {
				x1, x2 = x2, x1
			}

			hEdges = append(hEdges, hEdge{
				y:  a.y,
				x1: x1,
				x2: x2,
			})
		}
	}

	return
}

func getRedTiles(input [][]byte) []point {
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

func validCandidateRectangle(v1, v2 point, vertices []point, vEdges []vEdge, hEdges []hEdge) bool {

	x1 := min(v1.x, v2.x)
	x2 := max(v1.x, v2.x)
	y1 := min(v1.y, v2.y)
	y2 := max(v1.y, v2.y)

	if !pointInPolygon(point{x1, y1}, vertices) {
		return false
	}

	if !pointInPolygon(point{x1, y2}, vertices) {
		return false
	}

	if !pointInPolygon(point{x2, y1}, vertices) {
		return false
	}

	if !pointInPolygon(point{x2, y2}, vertices) {
		return false
	}

	for _, e := range vEdges {
		if rectHIntersectsPolyV(y1, x1, x2, e) {
			return false
		}

		if rectHIntersectsPolyV(y2, x1, x2, e) {
			return false
		}
	}

	for _, e := range hEdges {
		if rectVIntersectsPolyH(x1, y1, y2, e) {
			return false
		}

		if rectVIntersectsPolyH(x2, y1, y2, e) {
			return false
		}
	}

	return true
}

func pointInPolygon(p point, vertices []point) bool {
	inside := false
	px, py := p.x, p.y

	for i := range len(vertices) {
		a := vertices[i]
		b := vertices[(i+1)%len(vertices)]

		if a.x != b.x {
			continue
		}

		xEdge := a.x
		y1 := a.y
		y2 := b.y
		if y1 > y2 {
			y1, y2 = y2, y1
		}

		if py >= y1 && py < y2 && px < xEdge {
			inside = !inside
		}
	}

	return inside
}

func rectHIntersectsPolyV(yR, x1, x2 int, e vEdge) bool {
	if !(e.y1 < yR && yR < e.y2) {
		return false
	}

	return x1 < e.x && e.x < x2
}

func rectVIntersectsPolyH(xR, y1, y2 int, e hEdge) bool {
	if !(e.x1 < xR && xR < e.x2) {
		return false
	}

	return y1 < e.y && e.y < y2
}
