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

	startingNodeID := "you"
	targetNodeID := "out"

	sortedNodes := topologicalSort(nodes)
	paths := countPaths(startingNodeID, sortedNodes)

	elapsed := time.Since(start).Seconds()
	fmt.Printf("the total number of paths to 'out' is: %d\n", paths[targetNodeID])
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

	inPaths := make(map[string]int)
	nodeById := make(map[string]Node)

	for _, node := range nodes {
		nodeById[node.id] = node
		inPaths[node.id] = 0
	}

	for _, node := range nodes {
		for _, edge := range node.edges {
			inPaths[edge]++
		}
	}

	for _, node := range nodes {
		if inPaths[node.id] == 0 {
			queue = append(queue, node)
		}
	}

	for idx := 0; idx < len(queue); idx++ {
		currNode := queue[idx]
		result = append(result, currNode)

		for _, edge := range currNode.edges {
			inPaths[edge]--
			if inPaths[edge] == 0 {
				queue = append(queue, nodeById[edge])
			}
		}
	}

	fmt.Println("\t... complete")
	return result
}

func countPaths(startingID string, nodes []Node) map[string]int {
	fmt.Println("counting paths...")

	paths := make(map[string]int)
	paths[startingID] = 1

	for _, node := range nodes {
		for _, edge := range node.edges {
			paths[edge] += paths[node.id]
		}
	}

	fmt.Println("\t... complete")
	return paths
}
