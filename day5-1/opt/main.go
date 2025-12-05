package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Range struct {
	start int
	end   int
}

func main() {
	path := "../input.txt"
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	freshIngredientIDs := []Range{}
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}

		start, end, err := parseIDRange(line)
		if err != nil {
			panic(err)
		}

		newRange := Range{
			start: start,
			end:   end,
		}

		freshIngredientIDs = append(freshIngredientIDs, newRange)
	}

	mergedRanges := mergeRanges(freshIngredientIDs)

	totalFreshIngredients := 0
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}

		id, err := strconv.Atoi(line)
		if err != nil {
			panic(err)
		}

		if isFreshID(id, mergedRanges) {
			totalFreshIngredients++
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	fmt.Printf("the total number of fresh ingredients: %d\n", totalFreshIngredients)
}

func parseIDRange(input string) (start, end int, err error) {
	parts := strings.Split(input, "-")
	if len(parts) != 2 {
		err = fmt.Errorf("invalid ID range format: %s", input)
		return 0, 0, err
	}

	start, err = strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, err
	}

	end, err = strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, err
	}

	return start, end, nil
}

func mergeRanges(ranges []Range) []Range {
	if len(ranges) == 0 {
		return nil
	}

	sort.Slice(ranges, func(i, j int) bool {
		return ranges[i].start < ranges[j].start
	})

	merged := make([]Range, 0, len(ranges))

	current := ranges[0]

	for i := 1; i < len(ranges); i++ {
		r := ranges[i]

		if r.start <= current.end+1 {
			if r.end > current.end {
				current.end = r.end
			}
			continue
		}

		merged = append(merged, current)
		current = r
	}

	merged = append(merged, current)

	return merged
}

func isFreshID(id int, ranges []Range) bool {
	if len(ranges) == 0 {
		return false
	}

	lo, hi := 0, len(ranges)-1

	for lo <= hi {
		mid := (lo + hi) / 2
		r := ranges[mid]

		if id < r.start {
			hi = mid - 1
			continue
		}

		if id > r.end {
			lo = mid + 1
			continue
		}

		return true
	}

	return false
}
