package main

import "time"

type baseAlpha struct {
	start      time.Time
	end        time.Time
	current    time.Time
	currencies []string
	holdings   []float64
	tradeQueue []trade
}

func tradeOnTime() {

}
