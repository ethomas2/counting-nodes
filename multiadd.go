package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"sync"
)

var LIMIT uint64 = 20000000000
var ngoroutines = flag.Int("goroutines", 0, "number of goroutines to use")

func main() {
	cpuprofile := flag.String("cpuprofile", "", "write cpu profile to file")
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

	if flag.Arg(0) == "1" || flag.Arg(0) == "" {
		main1()
	} else if flag.Arg(0) == "2" {
		main2()
	} else {
		fmt.Println("\033[31;mBad input\033[m")
		os.Exit(1)
	}
}

func main1() {
	var total uint64
	var i uint64
	for i = 0; i < LIMIT; i++ {
		total += i
	}
	fmt.Println(total)
}

func main2() {
	var n uint64
	if *ngoroutines == 0 {
		n = uint64(runtime.GOMAXPROCS(0))
	} else {
		n = uint64(*ngoroutines)
	}
	fmt.Println("GOMAXPROCS", n)
	sums := make([]uint64, n)
	wg := sync.WaitGroup{}
	var i uint64
	for i = 0; i < n; i++ {
		wg.Add(1)
		go func(i uint64) {
			var start uint64 = (LIMIT / n) * i
			var end = start + (LIMIT / n)
			for j := start; j < end; j += 1 {
				sums[i] += j
			}
			wg.Done()
		}(uint64(i))
	}
	wg.Wait()

	var total uint64
	for _, s := range sums {
		total += s
	}
	fmt.Println(total)
}
