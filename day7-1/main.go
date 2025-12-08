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

	beamLocations := make(map[int]bool)
	beamLocations[sIndex] = true

	totalSplits := 0
	for scanner.Scan() {
		nextLine := scanner.Text()
		for v := range beamLocations {
			switch nextLine[v] {
			case '^':
				totalSplits++
				delete(beamLocations, v)
				beamLocations[v-1] = true
				beamLocations[v+1] = true
			}
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	fmt.Printf("the total number of splits: %d\n", totalSplits)
}
