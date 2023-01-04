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

func tradeOnTime() {
	//write alpha here
}
