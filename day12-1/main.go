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

	path := "./input.txt"
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	lines := bytes.Split(data, []byte("\n"))

	presents := []Present{}
	regions := []Region{}

	for i, line := range lines {
		if len(line) == 0 {
			continue
		}

		if line[1] == ':' {
			present := Present(lines[i+1 : i+4])
			presents = append(presents, present)
		}

		if len(line) < 3 {
			continue
		}

		if line[2] == 'x' {
			region, err := parseRegion(line)
			if err != nil {
				panic(err)
			}

			regions = append(regions, region)
		}

	}

	totalRegions := 0

	for _, region := range regions {
		requiredArea := 0
		for i, quantity := range region.Quantities {
			requiredArea += quantity * presents[i].Area()
		}

		if requiredArea <= region.Area() {
			totalRegions++
		}
	}

	elapsed := time.Since(start).Seconds()
	fmt.Printf("total runtime: %fs\n", elapsed)
	fmt.Printf("total number of regions that can fit the required presents: %d\n", totalRegions)
}

type Present [][]byte

func (p Present) Area() int {
	area := 0
	for _, row := range p {
		for _, c := range row {
			if c == '#' {
				area++
			}
		}
	}

	return area
}

type Region struct {
	X          int
	Y          int
	Quantities []int
}

func (r Region) Area() int {
	return r.X * r.Y
}

func parseRegion(line []byte) (Region, error) {
	var region Region

	tokens := bytes.Split(line, []byte(" "))

	dimensions := bytes.Split(tokens[0], []byte("x"))
	x, err := strconv.Atoi(string(dimensions[0]))
	if err != nil {
		return region, err
	}
	y, err := strconv.Atoi(string(bytes.TrimSuffix(dimensions[1], []byte(":"))))
	if err != nil {
		return region, err
	}

	region.X = x
	region.Y = y

	quantities := []int{}
	for _, token := range tokens[1:] {
		q, err := strconv.Atoi(string(token))
		if err != nil {
			return region, err
		}

		quantities = append(quantities, q)
	}

	if len(quantities) != 6 {
		return region, fmt.Errorf("unexpected quantities format: %s", line)
	}

	region.Quantities = quantities

	return region, nil
}
