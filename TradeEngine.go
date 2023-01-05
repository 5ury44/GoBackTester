package main

import (
	"fmt"
	"strings"
	"time"
)

type trade struct {
	currenciesKey int
	inverted      bool
	volume        float64
}

var mapCurrency map[string][]instant
var positions map[string]int
var first = ""

func initEngine(alpha baseAlpha) {
	initAlpha(alpha)
	findPairs(alpha.currencies)
	if alpha.downloadData {
		for _, i := range pairs {
			fmt.Println(i)
			whichCSV(alpha.start, alpha.end, currencyPairs[i])
		}
	}
	mapCurrency = currencyMap()
	positions = make(map[string]int)

	for key := range mapCurrency {
		positions[key] = binarySearch(mapCurrency[key], alpha.start)
	}
}

func executeEngine(alpha baseAlpha) {
	evalWorth(alpha)
	fmt.Println(" by start")
	for t := alpha.start; t.Before(alpha.end); t = t.Add(time.Millisecond * 1) {
		alpha.current = t
		for currencyExch := range positions {
			if mapCurrency[currencyExch][positions[currencyExch]].date.Before(t) {
				positions[currencyExch] += 1
			}
		}

		tradeOnTime(alpha)
		for _, tradeReg := range alpha.tradeQueue {
			currOpp := currencyPairs[tradeReg.currenciesKey]
			thisInstant :=
				mapCurrency[currOpp][positions[currOpp]]
			if tradeReg.inverted {
				alpha.holdings[currOpp[len(currencyPairs[tradeReg.currenciesKey])-3:]] -= tradeReg.volume
				alpha.holdings[currOpp[3:]] += tradeReg.volume * thisInstant.ask
			} else {
				alpha.holdings[currOpp[len(currOpp)-3:]] += tradeReg.volume / thisInstant.bid // todo check if bid is the correct term to use
				alpha.holdings[currOpp[3:]] -= tradeReg.volume
			}
		}
		alpha.tradeQueue = make([]trade, 0)
	}
	evalWorth(alpha)
	fmt.Println(" by end")
}

func binarySearch(instants []instant, target time.Time) int {
	low := 0
	high := len(instants) - 1

	for low <= high {
		mid := low + (high-low)/2
		midTime := instants[mid].date

		if midTime.Before(target) {
			low = mid + 1
		} else if midTime.After(target) {
			high = mid - 1
		} else {
			return mid
		}
	}

	if low < len(instants) && instants[low].date.After(target) {
		return low - 1
	}
	return low
}

func evalWorth(alpha baseAlpha) {
	transfer := make(map[string]float64, 0)
	for key, holding := range alpha.holdings {
		transfer[key] = holding
	}

	for key := range transfer {
		if first == "" || first == key {
			first = key
			continue
		}

		for _, s := range currencyPairs {
			if strings.Contains(s, key) && strings.Contains(s, first) {
				thisInstant := mapCurrency[s][positions[s]]
				if s[len(s)-3:] == key {
					transfer[first] += transfer[key] * thisInstant.ask
					transfer[key] = 0
				} else {
					transfer[first] += transfer[key] / thisInstant.bid //todo check if statement logic right
					transfer[key] = 0
				}
				break
			}
		}
	}
	fmt.Printf("%.10f", transfer[first])
	fmt.Print(" amount of " + first)
}
