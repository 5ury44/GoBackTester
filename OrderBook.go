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
	filesToParse := 2
	wg.Add(filesToParse)

	go func() {
		defer wg.Done()
		fmt.Println(len(processCSV("EURUSD-2022-11.csv")))
	}()

	go func() {
		defer wg.Done()
		fmt.Println(len(processCSV("EURUSD-2022-10.csv")))
	}()

	wg.Wait() // put everything to happen after below
	fmt.Println("loaded CSV files")
}
