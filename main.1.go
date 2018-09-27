package main

import (
	"fmt"
	"os"
	"strconv"
)

type tuple struct {
	x int
	y int
}

var N, M int

func getChildren(node tuple, visited map[tuple]struct{}) []tuple {
	var children []tuple
	deltas := []tuple{
		tuple{1, 0}, tuple{-1, 0}, tuple{0, 1}, tuple{0, -1},
	}
	for _, delta := range deltas {
		newNode := tuple{node.x + delta.x, node.y + delta.y}
		if newNode.x < 0 || newNode.x >= N || newNode.y < 0 || newNode.y >= M {
			continue
		}
		if _, ok := visited[newNode]; ok {
			continue

		}
		children = append(children, newNode)
	}
	return children
}

func dfs(node tuple, visited map[tuple]struct{}) (int, int) {
	children := getChildren(node, visited)
	if len(children) == 0 {
		return 1, 1
	}
	var totalPaths, totalTerminalPaths int = 0, 0
	for _, child := range children {
		visited[child] = struct{}{}
		nPaths, nTerminalPaths := dfs(child, visited)
		delete(visited, child)
		totalPaths += nPaths
		totalTerminalPaths += nTerminalPaths
	}
	return totalPaths + 1, totalTerminalPaths
}

func main() {
	mustConvert := func(s string) int {
		i, err := strconv.Atoi(s)
		if err != nil {
			panic(err)
		}
		return i
	}
	N, M = mustConvert(os.Args[1]), mustConvert(os.Args[2]) //globals.
	start := tuple{0, 0}
	fmt.Println(dfs(start, map[tuple]struct{}{start: struct{}{}}))
}
