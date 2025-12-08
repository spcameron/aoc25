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
	cache := distancesCache(sites)

	var edges []edge
	for pair, d := range cache {
		edges = append(edges, edge{pair[0], pair[1], d})
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

	um := newUnionMap(sites)

	for _, edge := range edges[0:1000] {
		um.union(edge.i, edge.j)
	}

	slices.SortFunc(um.sz, func(i, j int) int {
		return j - i
	})

	result := um.sz[0] * um.sz[1] * um.sz[2]

	fmt.Printf("the product of the three largest circuits: %d\n", result)
}

type site struct {
	x float64
	y float64
	z float64
}

func (s site) distance(t site) float64 {
	return math.Sqrt(math.Pow((s.x-t.x), 2) + math.Pow((s.y-t.y), 2) + math.Pow((s.z-t.z), 2))
}

type edge struct {
	i, j int
	dist float64
}

type unionMap struct {
	id    []int
	sz    []int
	count int
}

func newUnionMap[T any](sites []T) *unionMap {
	n := len(sites)

	um := unionMap{}
	um.id = make([]int, n)
	um.sz = make([]int, n)
	um.count = n

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

func (um *unionMap) connected(p, q int) bool {
	return um.find(p) == um.find(q)
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

func sitesFromRaw(data []byte) []site {
	fmt.Println("building sites")
	lines := bytes.Split(data, []byte("\n"))
	sites := make([]site, len(lines))
	for i, line := range lines {
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

		sites[i] = site{float64(x), float64(y), float64(z)}
	}

	fmt.Println("sites built")
	return sites
}

func distancesCache(data []site) map[[2]int]float64 {
	fmt.Println("building cache")
	cache := make(map[[2]int]float64)
	for i := 0; i < len(data); i++ {
		for j := i + 1; j < len(data); j++ {
			cache[[2]int{i, j}] = data[i].distance(data[j])
		}
	}

	fmt.Println("cache built")
	return cache
}
