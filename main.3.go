package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"strconv"
	"sync/atomic"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var ngoroutines = flag.Int("goroutines", 0, "number of goroutines to use")

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

func sum(arr []paddedInt) int {
	var ret int
	for _, x := range arr {
		ret += x.n
	}
	return ret
}

type paddedInt struct {
	n       int
	padding [64]byte
}

func solve(gNode, nRows, nCols uint64) (int, int) {
	var NUMGOROUTINES int
	if *ngoroutines == 0 {
		NUMGOROUTINES = runtime.GOMAXPROCS(0)
	} else {
		NUMGOROUTINES = *ngoroutines
	}
	nLeafNodes := make([]paddedInt, NUMGOROUTINES)
	nNodes := make([]paddedInt, NUMGOROUTINES)
	result := make(chan tuple, 1)
	var outstanding uint64 = 1

	nodeschan := make(chan TNode, NUMGOROUTINES)
	for i := 0; i < NUMGOROUTINES; i++ {
		ii := i
		visit := func(tnode TNode, children []TNode) {
			atomic.AddUint64(&outstanding, uint64(len(children)))
			if len(children) == 0 {
				nLeafNodes[ii].n += 1
			}
			nNodes[ii].n += 1
			if newval := atomic.AddUint64(&outstanding, ^uint64(0)); newval == 0 {
				result <- tuple{sum(nNodes), sum(nLeafNodes)}
			}
		}
		go traverse(nodeschan, visit, getChildren)
	}
	tnode := TNode{gNode, 1 << gNode}
	nodeschan <- tnode

	tup := <-result
	return tup.fst, tup.snd
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

	nRows, nCols = uint64(mustAToi(flag.Arg(0))), uint64(mustAToi(flag.Arg(1)))
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
