package main

import (
	"bytes"
	"fmt"
	"os"
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
	nodes := parseLines(lines)

	startingNodeID := "svr"
	targetNodeID := "out"

	sortedNodes := topologicalSort(nodes)
	paths := countPaths(startingNodeID, sortedNodes)

	elapsed := time.Since(start).Seconds()
	fmt.Printf("the total number of paths to 'out' with waypoints 'dac' & 'fft' is: %d\n", paths[targetNodeID][3])
	fmt.Printf("total time: %vs", elapsed)
}

type Node struct {
	id    string
	edges []string
}

func parseLines(lines [][]byte) []Node {
	fmt.Println("parsing lines...")

	nodes := []Node{}
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}

		node := parseLine(line)
		nodes = append(nodes, node)
	}

	fmt.Println("\t... complete")
	return nodes
}

func parseLine(line []byte) Node {
	tokens := bytes.Split(line, []byte(" "))

	id := string(bytes.TrimSuffix(tokens[0], []byte(":")))

	edges := []string{}
	for _, token := range tokens[1:] {
		edges = append(edges, string(token))
	}

	return Node{id, edges}
}

func topologicalSort(nodes []Node) []Node {
	fmt.Println("sorting nodes...")

	queue := []Node{}
	result := []Node{}

	inDegree := make(map[string]int)
	nodeByID := make(map[string]Node)

	for _, node := range nodes {
		nodeByID[node.id] = node
		inDegree[node.id] = 0
	}

	for _, node := range nodes {
		for _, edge := range node.edges {
			inDegree[edge]++
		}
	}

	for _, node := range nodes {
		if inDegree[node.id] == 0 {
			queue = append(queue, node)
		}
	}

	for idx := 0; idx < len(queue); idx++ {
		currNode := queue[idx]
		result = append(result, currNode)

		for _, edge := range currNode.edges {
			inDegree[edge]--
			if inDegree[edge] == 0 {
				queue = append(queue, nodeByID[edge])
			}
		}
	}

	fmt.Println("\t... complete")
	return result
}

func countPaths(startingID string, nodes []Node) map[string][4]int {
	fmt.Println("counting paths...")

	paths := map[string][4]int{}
	paths[startingID] = [4]int{1, 0, 0, 0}

	for _, node := range nodes {
		for _, edge := range node.edges {
			for S := range 4 {
				count := paths[node.id][S]
				if count == 0 {
					continue
				}

				newState := S
				if edge == "dac" {
					newState |= 1
				}

				if edge == "fft" {
					newState |= 2
				}

				arr := paths[edge]
				arr[newState] += count
				paths[edge] = arr
			}
		}
	}

	fmt.Println("\t... complete")
	return paths
}
