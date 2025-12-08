package main

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"slices"
	"strconv"
)

func main() {
	path := "./input.txt"
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	sites := sitesFromRaw(data)
	edges := makeEdgesSortedByDistance(sites)
	um := newUnionMap(sites)

	lastEdge := edge{}
	for _, edge := range edges {
		um.union(edge.i, edge.j)
		if um.count == 1 {
			lastEdge = edge
			break
		}
	}

	result := sites[lastEdge.i].x * sites[lastEdge.j].x
	fmt.Printf("product of x-coords of last two junction boxes: %d\n", result)
}

type site struct {
	x int
	y int
	z int
}

func (s site) distance(t site) float64 {
	return math.Pow(float64(s.x-t.x), 2) + math.Pow(float64(s.y-t.y), 2) + math.Pow(float64(s.z-t.z), 2)
}

func sitesFromRaw(data []byte) []site {
	lines := bytes.Split(data, []byte("\n"))
	sites := []site{}
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}

		parts := bytes.Split(line, []byte(","))
		if len(parts) != 3 {
			panic(fmt.Sprintf("unexpected line format: %s", string(line)))
		}

		x, err := strconv.Atoi(string(parts[0]))
		if err != nil {
			panic(err)
		}

		y, err := strconv.Atoi(string(parts[1]))
		if err != nil {
			panic(err)
		}

		z, err := strconv.Atoi(string(parts[2]))
		if err != nil {
			panic(err)
		}

		sites = append(sites, site{x, y, z})
	}

	return sites
}

type edge struct {
	i    int
	j    int
	dist float64
}

func makeEdgesSortedByDistance(sites []site) []edge {
	n := len(sites)
	edges := make([]edge, n*(n-1)/2)
	idx := 0
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			edges[idx] = edge{i, j, sites[i].distance(sites[j])}
			idx++
		}
	}

	slices.SortFunc(edges, func(a, b edge) int {
		switch {
		case a.dist < b.dist:
			return -1
		case a.dist > b.dist:
			return 1
		default:
			return 0
		}
	})

	return edges
}

type unionMap struct {
	id    []int
	sz    []int
	count int
}

func newUnionMap[T any](sites []T) *unionMap {
	n := len(sites)
	um := unionMap{
		id:    make([]int, n),
		sz:    make([]int, n),
		count: n,
	}

	for i := range n {
		um.id[i] = i
		um.sz[i] = 1
	}

	return &um
}

func (um *unionMap) find(p int) int {
	for p != um.id[p] {
		p = um.id[p]
	}

	return p
}

func (um *unionMap) union(p, q int) {
	i := um.find(p)
	j := um.find(q)
	if i == j {
		return
	}

	if um.sz[i] < um.sz[j] {
		um.id[i] = j
		um.sz[j] += um.sz[i]
	} else {
		um.id[j] = i
		um.sz[i] += um.sz[j]
	}

	um.count--
}
