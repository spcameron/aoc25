package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	path := "../input.txt"
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

	width := len(firstLine)
	beamsCurrent := make([]int64, width)
	beamsCurrent[sIndex] = 1

	for scanner.Scan() {
		nextLine := scanner.Text()
		beamsNext := make([]int64, width)

		for c := range width {
			v := beamsCurrent[c]

			if v == 0 {
				continue
			}

			switch nextLine[c] {
			case '.':
				beamsNext[c] += v
			case '^':
				if c-1 >= 0 {
					beamsNext[c-1] += v
				}
				if c+1 < width {
					beamsNext[c+1] += v
				}
			}
		}

		beamsCurrent = beamsNext
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	var totalTimelines int64
	for _, v := range beamsCurrent {
		totalTimelines += v
	}

	fmt.Printf("the total number of timelines: %d\n", totalTimelines)
}
