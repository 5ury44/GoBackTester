package main

var queue = meanQueue{
	size:  500,
	items: make([]float64, 0),
}
var enter = float64(0)
var foundArc = false

func exampleReversion(alpha baseAlpha, positions map[string]int) baseAlpha {
	currentAsk := mapCurrency["eurusd"][positions["eurusd"]].ask
	queue.Enqueue(currentAsk)

	if queue.Len() < 500 {
		return alpha
	}

	if foundArc == false && enter == 0 && queue.predictNextPoint() > currentAsk && queue.PosSlope() {
		enter = currentAsk
		foundArc = true
		return makeTrade(alpha, "usd", "eur", alpha.holdings["usd"])
	} else if foundArc == true && queue.predictNextPoint() < currentAsk && queue.NegSlope() && enter < currentAsk {
		enter = 0
		foundArc = false
		return makeTrade(alpha, "eur", "usd", alpha.holdings["eur"])
	}

	return alpha

}
