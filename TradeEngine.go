package main

import (
	"fmt"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
	"os"
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
var history = make([]opts.LineData, 0)
var initialValue float64 = 0

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

func executeEngine(alpha baseAlpha, resolution int) {
	initialValue = evalWorth(alpha)
	fmt.Printf("%.10f", initialValue)
	fmt.Println(" amount of " + first + " portfolio worth at the start")
	count := 0
	for t := alpha.start; t.Before(alpha.end); t = t.Add(time.Millisecond * 1) {

		count++
		alpha.current = t
		for currencyExch := range positions {
			if mapCurrency[currencyExch][positions[currencyExch]].date.Before(t) {
				positions[currencyExch] += 1
			}
		}
		if count-1 == 0 || (count-1)%100 == 0 {
			ba := findMaxAskIndexWithin1000ms(positions["eurusd"])
			newpos := make(map[string]int, 0)
			newpos["eurusd"] = ba[0]
			alpha := tradeOnTime(alpha, newpos)

			for _, tradeReg := range alpha.tradeQueue {
				currOpp := currencyPairs[tradeReg.currenciesKey]
				thisInstant :=
					mapCurrency[currOpp][newpos["eurusd"]]
				if tradeReg.inverted {
					alpha.holdings[currOpp[len(currencyPairs[tradeReg.currenciesKey])-3:]] -= tradeReg.volume
					alpha.holdings[currOpp[3:]] += tradeReg.volume * thisInstant.ask
				} else {
					alpha.holdings[currOpp[len(currOpp)-3:]] += tradeReg.volume / mapCurrency[currOpp][ba[1]].bid // todo check if bid is the correct term to use
					alpha.holdings[currOpp[3:]] -= tradeReg.volume
				}
			}

			alpha.tradeQueue = make([]trade, 0)
		}
		graphInit(resolution, count, alpha)
	}
	fmt.Printf("%.10f", evalWorth(alpha))
	fmt.Println(" amount of " + first + " portfolio worth by the end")
	createGraph()
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

func evalWorth(alpha baseAlpha) float64 {
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
					transfer[first] += transfer[key] / thisInstant.bid //todo maybe standardize out currency?
					transfer[key] = 0
				}
				break
			}
		}
	}
	return transfer[first]
}

func graphInit(resolution int, current int, alpha baseAlpha) {
	if resolution > 0 && current%resolution == 0 {
		history = append(history, opts.LineData{Value: evalWorth(alpha) / initialValue})
	}
}

func createGraph() {
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeWesteros}),
		charts.WithTitleOpts(opts.Title{
			Title:    "Percent Change in Portfolio Worth Over Given Period",
			Subtitle: "By GoBackTester Forex Backtester",
		}))
	charts.WithYAxisOpts(opts.YAxis{Name: "percent change"})
	line.AddSeries("Percent Growth", history).SetXAxis(len(history))
	f, _ := os.Create("line" + time.Now().String() + ".html")
	line.Render(f)
}

func findMaxAskIndexWithin1000ms(pos int) []int {
	instants := mapCurrency["eurusd"]
	maxAsk := -1.0
	maxAskIndex := -1
	maxBid := -1.0
	maxBidIndex := -1
	endTime := instants[pos].date.Add(100 * time.Millisecond)
	for i := pos; i < len(instants) && instants[i].date.Before(endTime); i++ {
		if instants[i].ask > maxAsk {
			maxAsk = instants[i].ask
			maxAskIndex = i
		}
		if 1/instants[i].bid > maxBid {
			maxBid = instants[i].bid
			maxBidIndex = i
		}
	}
	return []int{maxAskIndex, maxBidIndex}
}
