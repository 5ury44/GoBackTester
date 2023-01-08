package main

type meanQueue struct {
	size  int
	items []float64
}

var total = float64(0)

func (q *meanQueue) Enqueue(item float64) {
	if len(q.items) == q.size {
		q.items = q.items[1:]
	}
	total += item
	q.items = append(q.items, item)
}

func (q *meanQueue) Dequeue() float64 {
	item := q.items[0]
	q.items = q.items[1:]
	total -= item
	return item
}

func (q *meanQueue) Len() int {
	return len(q.items)
}

func (q *meanQueue) Mean() float64 {
	return total / float64(q.Len())
}

func (q *meanQueue) PosSlope() bool {
	sum := 0.0
	for i := 1; i < 5; i++ {
		sum += q.items[len(q.items)-i] - q.items[len(q.items)-i-1]
	}
	return sum/5 > 0
}

func (q *meanQueue) NegSlope() bool {
	sum := 0.0
	for i := 1; i < 5; i++ {
		sum += q.items[len(q.items)-i] - q.items[len(q.items)-i-1]
	}
	return sum/5 < 0
}

func (q *meanQueue) predictNextPoint() float64 {
	points := q.items

	// Calculate the average difference between each point
	var sum float64
	for i := 1; i < len(points); i++ {
		sum += points[i] - points[i-1]
	}
	averageDifference := sum / float64(len(points)-1)

	// Return the next predicted point by adding the average difference to the last point
	return points[len(points)-1] + averageDifference
}
