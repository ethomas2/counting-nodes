package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"strconv"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

type TNode struct {
	gNode   uint64
	bitmask uint64
}

type tuple struct {
	fst int
	snd int
}

func dfsTree(root TNode, getChildren func(TNode) []TNode) (int, int) {
	children := getChildren(root)
	if len(children) == 0 {
		return 1, 1
	}

	nNodes, nLeaves := 0, 0
	for _, child := range children {
		x, y := dfsTree(child, getChildren)
		nNodes += x
		nLeaves += y
	}
	return nNodes + 1, nLeaves
}

func solve(gNode, nRows, nCols uint64) (int, int) {
	getChildren := func(tnode TNode) []TNode {
		deltas := []tuple{
			tuple{1, 0},
			tuple{-1, 0},
			tuple{0, 1},
			tuple{0, -1},
		}
		children := make([]TNode, 0, 4)
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
	tnode := TNode{gNode, 1 << gNode}
	return dfsTree(tnode, getChildren)
}

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer func() {
			pprof.StopCPUProfile()
			fmt.Println("stopped cpu profile")
		}()
	}

	nRows, nCols := mustAToi(flag.Arg(0)), mustAToi(flag.Arg(1))
	// TODO: allow custom start nodes
	// TODO: make gNodes a type (type gNode int)
	fmt.Println(solve(0, uint64(nRows), uint64(nCols)))

}

func mustAToi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i

}
