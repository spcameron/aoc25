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
	edges := getEdges(vertices)

	maxArea := 0
	for i := 0; i < len(vertices); i++ {
		for j := 0; j < len(vertices); j++ {
			r := Rectangle{vertices[i], vertices[j]}
			if r.Area() < maxArea {
				continue
			}

			if !r.ContainedWithin(edges) {
				continue
			}

			maxArea = r.Area()
		}
	}

	elapsed := time.Since(start).Seconds()
	fmt.Printf("the largest area of any valid rectangle is: %d\n", maxArea)
	fmt.Printf("time elapsed: %vs\n", elapsed)
}

type Point struct {
	x int
	y int
}

type Edge struct {
	v1 Point
	v2 Point
}

func (e Edge) isHorizontal() bool {
	return e.v1.y == e.v2.y
}

func (e Edge) isVertical() bool {
	return e.v1.x == e.v2.x
}

func NormalizeEdge(e Edge) Edge {
	a := e.v1
	b := e.v2

	if a.x > b.x || (a.x == b.x && a.y > b.y) {
		a, b = b, a
	}

	return Edge{a, b}
}

type Rectangle struct {
	corner1 Point
	corner2 Point
}

func (r Rectangle) Area() int {
	dx := r.corner1.x - r.corner2.x
	if dx < 0 {
		dx = -dx
	}

	dy := r.corner1.y - r.corner2.y
	if dy < 0 {
		dy = -dy
	}

	return (dx + 1) * (dy + 1)
}

func (r Rectangle) ContainedWithin(edges []Edge) bool {
	if !r.CornersContained(edges) {
		return false
	}

	if r.IntersectedByPolygon(edges) {
		return false
	}

	return true
}

func (r Rectangle) CornersContained(edges []Edge) bool {
	if !pointInPolygon(Point{r.corner1.x, r.corner1.y}, edges) {
		return false
	}

	if !pointInPolygon(Point{r.corner1.x, r.corner2.y}, edges) {
		return false
	}

	if !pointInPolygon(Point{r.corner2.x, r.corner1.y}, edges) {
		return false
	}

	if !pointInPolygon(Point{r.corner2.x, r.corner2.y}, edges) {
		return false
	}

	return true
}

func (r Rectangle) IntersectedByPolygon(edges []Edge) bool {
	for _, e := range edges {
		if segmentsIntersect(e, r.TopEdge()) {
			return true
		}

		if segmentsIntersect(e, r.BottomEdge()) {
			return true
		}

		if segmentsIntersect(e, r.RightEdge()) {
			return true
		}

		if segmentsIntersect(e, r.LeftEdge()) {
			return true
		}
	}

	return false
}

func (r Rectangle) TopEdge() Edge {
	y := min(r.corner1.y, r.corner2.y)
	x1, x2 := r.corner1.x, r.corner2.x
	if x1 > x2 {
		x1, x2 = x2, x1
	}

	return Edge{Point{x1, y}, Point{x2, y}}
}

func (r Rectangle) BottomEdge() Edge {
	y := max(r.corner1.y, r.corner2.y)
	x1, x2 := r.corner1.x, r.corner2.x
	if x1 > x2 {
		x1, x2 = x2, x1
	}

	return Edge{Point{x1, y}, Point{x2, y}}
}

func (r Rectangle) LeftEdge() Edge {
	x := min(r.corner1.x, r.corner2.x)
	y1, y2 := r.corner1.y, r.corner2.y
	if y1 > y2 {
		y1, y2 = y2, y1
	}

	return Edge{Point{x, y1}, Point{x, y2}}
}

func (r Rectangle) RightEdge() Edge {
	x := max(r.corner1.x, r.corner2.x)
	y1, y2 := r.corner1.y, r.corner2.y
	if y1 > y2 {
		y1, y2 = y2, y1
	}

	return Edge{Point{x, y1}, Point{x, y2}}
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

func getEdges(vertices []Point) []Edge {
	edges := []Edge{}
	for i := range len(vertices) {
		a := vertices[i]
		b := vertices[(i+1)%len(vertices)]

		edges = append(edges, NormalizeEdge(Edge{a, b}))
	}

	return edges
}

func pointInPolygon(p Point, edges []Edge) bool {
	inside := false

	for _, e := range edges {
		if !e.isVertical() {
			continue
		}

		xEdge := e.v1.x
		y1, y2 := e.v1.y, e.v2.y
		if y1 > y2 {
			y1, y2 = y2, y1
		}

		if p.y >= y1 && p.y < y2 && p.x < xEdge {
			inside = !inside
		}
	}

	return inside
}

func segmentsIntersect(p, q Edge) bool {
	o1 := orient(p.v1, p.v2, q.v1)
	o2 := orient(p.v1, p.v2, q.v2)
	o3 := orient(q.v1, q.v2, p.v1)
	o4 := orient(q.v1, q.v2, p.v2)

	if oppSign(o1, o2) && oppSign(o3, o4) {
		return true
	}

	if o1 == 0 && onSegment(p.v1, p.v2, q.v1) {
		return true
	}

	if o2 == 0 && onSegment(p.v1, p.v2, q.v2) {
		return true

	}

	if o3 == 0 && onSegment(q.v1, q.v2, p.v1) {
		return true

	}

	if o4 == 0 && onSegment(q.v1, q.v2, p.v2) {
		return true

	}

	return false
}

func orient(a, b, c Point) int {
	return (b.x-a.x)*(c.y-a.y) - (b.y-a.y)*(c.x-a.x)
}

func onSegment(a, b, c Point) bool {
	if !(min(a.x, b.x) < c.x && c.x < max(a.x, b.x)) {
		return false
	}

	if !(min(a.y, b.y) < c.y && c.y < max(a.y, b.y)) {
		return false
	}

	return true
}

func oppSign(a, b int) bool {
	return (a < 0 && b > 0) || (a > 0 && b < 0)
}
