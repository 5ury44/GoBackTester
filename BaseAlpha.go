package main

import "time"

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
	executeEngine(alpha)
}

func initAlpha(alpha baseAlpha) {
	alpha.current = alpha.start
	// write alpha initialization code here if needed
}

func tradeOnTime(alpha baseAlpha) {
	//write alpha here
}
