package main

import (
	"runtime"
	"sync"
	"time"
)

var wg sync.WaitGroup

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	/*filesToParse := 2
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
	fmt.Println("loaded CSV files")*/

	/*current := time.Date(2022, 11, 1, 1, 1, 1, 0, time.UTC)
	next := time.Date(2022, 11, 3, 1, 1, 1, 0, time.UTC)
	whichCSV(current, next, "eurusd")*/

	/*testMap := currencyMap([]string{"usd", "eur"})
	fmt.Println(len(testMap["eurusd"]))*/

	alpha := baseAlpha{
		start:      time.Date(2022, 11, 1, 1, 1, 1, 0, time.UTC),
		end:        time.Date(2022, 11, 3, 1, 1, 1, 0, time.UTC),
		current:    time.Date(2022, 11, 1, 1, 1, 1, 0, time.UTC), //set to start by default
		currencies: []string{"eur", "usd"},
		holdings: map[string]float64{
			"eur": 0.5,
			"usd": 0.5,
		},
		tradeQueue:   make([]trade, 1),
		downloadData: false,
	}
	initEngine(alpha)
	executeEngine(alpha)
	wg.Wait()
}
