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
	path := "./input.txt"
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	freshIngredientIDs := []Range{}
	for scanner.Scan() {
		line := scanner.Text()
		// fmt.Printf("processing range line: %s\n", line)
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

	sort.Slice(freshIngredientIDs, func(i, j int) bool {
		return freshIngredientIDs[i].start < freshIngredientIDs[j].start
	})

	totalFreshIngredients := 0
	for scanner.Scan() {
		line := scanner.Text()
		// fmt.Printf("process ID line: %s\n", line)
		if line == "" {
			break
		}

		id, err := strconv.Atoi(line)
		if err != nil {
			panic(err)
		}

		if isFreshID(id, freshIngredientIDs) {
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

func isFreshID(id int, idRanges []Range) bool {
	for _, r := range idRanges {
		if id >= r.start && id <= r.end {
			return true
		}

		if id < r.start {
			break
		}
	}

	return false
}
