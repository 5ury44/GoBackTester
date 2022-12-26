package main

import (
	"fmt"
	"runtime"
	"sync"
)

var wg sync.WaitGroup

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	const filesToParse = 2
	wg.Add(filesToParse)

	go fmt.Print(len(processCSV("EURUSD-2022-11.csv")))
	go fmt.Print(len(processCSV("EURUSD-2022-11.csv")))

	wg.Wait() // put everything to happen after below
}
