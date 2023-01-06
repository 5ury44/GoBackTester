package main

import (
	"strings"
	"time"
)

type baseAlpha struct {
	start        time.Time
	end          time.Time
	current      time.Time
	currencies   []string
	holdings     map[string]float64
	tradeQueue   []trade
	downloadData bool
}

func newAlpha(start time.Time, end time.Time, curr []string, holdings map[string]float64, download bool) {
	alpha := baseAlpha{
		start:        start,
		end:          end,
		current:      start, //set to start by default
		currencies:   curr,
		holdings:     holdings,
		tradeQueue:   make([]trade, 0),
		downloadData: download,
	}
	initEngine(alpha)
	executeEngine(alpha, 60000) // higher resolution to run faster but worse graph quality (in ms intervals)
}

func initAlpha(alpha baseAlpha) {
	alpha.current = alpha.start
	// write alpha initialization code here if needed
}

func tradeOnTime(alpha baseAlpha) {
	//write alpha here
	//use make trade to make a trade which will be added to the queue
}

func makeTrade(alpha baseAlpha, currency1 string, currency2 string, volume float64) {
	for i, s := range currencyPairs {
		if strings.Contains(s, currency1) && strings.Contains(s, currency2) {
			invert := s[len(s)-3:] == currency1
			alpha.tradeQueue = append(alpha.tradeQueue, trade{i, invert, volume})
			break
		}
	}
}
