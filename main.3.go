package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"sync/atomic"
)

var nRows, nCols uint64

type TNode struct {
	gNode   uint64
	bitmask uint64
}

type tuple struct {
	fst int
	snd int
}

func getChildren(tnode TNode) []TNode {
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

func solve(gNode, nRows, nCols uint64) (uint64, uint64) {
	var nLeafNodes, nNodes, outstanding uint64
	type tuple struct {
		fst uint64
		snd uint64
	}

	result := make(chan tuple, 1)
	outstanding = 1
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
	if len(os.Args) < 3 {
		fmt.Println("\033[31mGimme some ints yo!\033[m")
		os.Exit(1)
	}

	nRows, nCols = uint64(mustAToi(os.Args[1])), uint64(mustAToi(os.Args[2]))
	// TODO: allow custom start nodes
	// TODO: make gNodes a type (type gNode int)
	fmt.Println(solve(0, uint64(nRows), uint64(nCols)))

}

func mustAToi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		fmt.Println("\033[31mThose don't look like ints\033[m")
		os.Exit(1)
	}
	return i

}
