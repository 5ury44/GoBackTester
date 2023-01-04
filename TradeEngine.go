package main

import "time"

type trade struct {
	currenciesKey int
	inverted      bool
	volume        float64
}

var mapCurrency map[string][]instant
var positions map[string]int

func initEngine(alpha baseAlpha) {
	findPairs(alpha.currencies)
	if alpha.downloadData {
		for i := range pairs {
			whichCSV(alpha.start, alpha.end, currencyPairs[i])
		}
	}
	mapCurrency = currencyMap()

	for key := range mapCurrency {
		positions[key] = binarySearch(mapCurrency[key], alpha.start)
	}
}

func executeEngine(alpha baseAlpha) {
	for t := alpha.start; t.Before(alpha.end); t = t.Add(time.Millisecond * 1) {
		alpha.current = t
		for currencyExch := range positions {
			if mapCurrency[currencyExch][positions[currencyExch]].date.Before(t) {
				positions[currencyExch] += 1
			}
		}

		tradeOnTime()
		for _, tradeReg := range alpha.tradeQueue {
			thisInstant := mapCurrency[currencyPairs[tradeReg.currenciesKey]][positions[currencyPairs[tradeReg.currenciesKey]]]
			if tradeReg.inverted {
				alpha.holdings[currencyPairs[tradeReg.currenciesKey][len(currencyPairs[tradeReg.currenciesKey])-3:]] -=
					tradeReg.volume
				alpha.holdings[currencyPairs[tradeReg.currenciesKey][3:]] += tradeReg.volume * thisInstant.ask
			} else {
				alpha.holdings[currencyPairs[tradeReg.currenciesKey][len(currencyPairs[tradeReg.currenciesKey])-3:]] +=
					tradeReg.volume / thisInstant.bid // todo check if bid is the correct term to use
				alpha.holdings[currencyPairs[tradeReg.currenciesKey][3:]] -= tradeReg.volume
			}
		}
	}
}

func binarySearch(instants []instant, target time.Time) int {
	low := 0
	high := len(instants) - 1

	for low <= high {
		mid := low + (high-low)/2

		if target.Before(instants[mid].date) {
			high = mid - 1
		} else if target.After(instants[mid].date) {
			low = mid + 1
		} else {
			return mid
		}
	}

	if low >= len(instants) {
		return len(instants) - 1
	} else if high < 0 {
		return 0
	} else if target.Sub(instants[low].date) < instants[high].date.Sub(target) {
		return low
	} else {
		return high
	}
}
