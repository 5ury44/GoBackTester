package main

import "time"

type trade struct {
	currenciesKey int
	inverted      bool
	volume        float64
}

var mapCurrency map[string][]instant

func initEngine(alpha baseAlpha) {
	findPairs(alpha.currencies)
	for i := range pairs {
		whichCSV(alpha.start, alpha.end, currencyPairs[i])
	}
	mapCurrency = currencyMap()

}

func executeEngine(alpha baseAlpha) {

}

func binarySearch(instants []instant, t time.Time) int {
	low := 0
	high := len(instants) - 1

	for low <= high {
		mid := low + (high-low)/2
		midTime := instants[mid].date

		if midTime.Equal(t) || midTime.Before(t) {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}

	return high
}
