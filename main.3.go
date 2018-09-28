package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"sync/atomic"
)

type TNode struct {
	gNode   uint64
	bitmask uint64
}

type tuple struct {
	fst int
	snd int
}

// func dfsTree(root TNode, getChildren func(TNode) []TNode) (int, int) {
// 	children := getChildren(root)
// 	if len(children) == 0 {
// 		return 1, 1
// 	}

// 	nNodes, nLeaves := 0, 0
// 	for _, child := range children {
// 		x, y := dfsTree(child, getChildren)
// 		nNodes += x
// 		nLeaves += y
// 	}
// 	return nNodes + 1, nLeaves
// }

func solve(gNode, nRows, nCols uint64) (uint64, uint64) {
	getChildren := func(tnode TNode) []TNode {
		deltas := []tuple{
			tuple{1, 0},
			tuple{-1, 0},
			tuple{0, 1},
			tuple{0, -1},
		}
		var children []TNode
		for _, delta := range deltas {
			var deltar, deltac int = delta.fst, delta.snd
			var r, c uint64 = tnode.gNode / nCols, tnode.gNode % nCols
			ooBounds := ((r <= 0 && deltar == -1) ||
				(c <= 0 && deltac == -1) ||
				(r >= nRows-1 && deltar == 1) ||
				(c >= nCols-1 && deltac == 1))
			if ooBounds {
				continue
			}
			r += uint64(deltar)
			c += uint64(deltac)

			newGNode := r*nCols + c
			visited := (1<<newGNode)&tnode.bitmask != 0
			if visited {
				continue
			}
			children = append(children, TNode{gNode: newGNode, bitmask: tnode.bitmask | (1 << newGNode)})
		}
		return children
	}

	var nLeafNodes, nNodes, outstanding uint64
	outstanding = 1
	type tuple struct {
		fst uint64
		snd uint64
	}
	result := make(chan tuple, 1)
	visit := func(tnode TNode, children []TNode) {
		atomic.AddUint64(&outstanding, uint64(len(children)))
		if len(children) == 0 {
			atomic.AddUint64(&nLeafNodes, 1)
		}
		atomic.AddUint64(&nNodes, 1)
		if newval := atomic.AddUint64(&outstanding, ^uint64(0)); newval == 0 {
			fmt.Println("DONE !!!", nLeafNodes, nNodes)
			result <- tuple{nNodes, nLeafNodes}
		}
	}
	NUMGOROUTINES := runtime.NumCPU()
	nodeschan := make(chan TNode, NUMGOROUTINES)
	for i := 0; i < NUMGOROUTINES; i++ {
		go traverse(nodeschan, visit, getChildren)
	}
	tnode := TNode{gNode, 1 << gNode}
	nodeschan <- tnode

	tup := <-result
	return tup.fst, tup.snd
}

func main() {
	nRows, nCols := mustAToi(os.Args[1]), mustAToi(os.Args[2])
	// TODO: allow custom start nodes
	// TODO: make gNodes a type (type gNode int)
	fmt.Println(solve(0, uint64(nRows), uint64(nCols)))

}

func mustAToi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		fmt.Println("\033[31mGimme some ints yo!\033[m")
		os.Exit(1)
	}
	return i

}
