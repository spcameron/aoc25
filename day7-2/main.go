package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	path := "./input.txt"
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Scan()
	firstLine := scanner.Text()

	sIndex := -1
	for i, c := range firstLine {
		if c == 'S' {
			sIndex = i
			break
		}
	}

	if sIndex == -1 {
		panic("did not find S index")
	}

	beamsToTimelines := map[int]int{sIndex: 1}
	for scanner.Scan() {
		nextLine := scanner.Text()
		for k, v := range beamsToTimelines {
			switch nextLine[k] {
			case '.':
				continue
			case '^':
				delete(beamsToTimelines, k)
				beamsToTimelines[k-1] += v
				beamsToTimelines[k+1] += v
			}
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	totalTimelines := 0
	for _, v := range beamsToTimelines {
		totalTimelines += v
	}

	fmt.Printf("the total number of timelines: %d\n", totalTimelines)
}
