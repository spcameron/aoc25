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
	beamCurrent := make([]bool, width)
	beamCurrent[sIndex] = true

	totalSplits := 0
	for scanner.Scan() {
		nextLine := scanner.Text()
		beamNext := make([]bool, width)

		for c := range width {
			if !beamCurrent[c] {
				continue
			}

			switch nextLine[c] {
			case '^':
				totalSplits++
				if c-1 >= 0 {
					beamNext[c-1] = true
				}
				if c+1 < width {
					beamNext[c+1] = true
				}
			default:
				beamNext[c] = true
			}
		}

		beamCurrent = beamNext
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	fmt.Printf("the total numberof splits: %d\n", totalSplits)
}
