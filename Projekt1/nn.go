package main

import (
	"math"
	"time"
)

// SolveNN implementuje prosty algorytm Nearest Neighbor (bez rozgałęzień).
// Zawsze wybiera pierwsze napotkane miasto o najmniejszym koszcie.
func (t TSPInstance) SolveNN(startCity int) Result {
	start := time.Now()
	n := t.Size

	visited := make([]bool, n)
	visited[startCity] = true

	path := []int{startCity}
	totalCost := 0
	currentCity := startCity

	for step := 1; step < n; step++ {
		minEdgeCost := math.MaxInt32
		nextCity := -1

		for j := 0; j < n; j++ {
			if !visited[j] {
				cost := t.Matrix[currentCity][j]
				if cost < minEdgeCost {
					minEdgeCost = cost
					nextCity = j
				}
			}
		}

		visited[nextCity] = true
		path = append(path, nextCity)
		totalCost += minEdgeCost
		currentCity = nextCity
	}

	// Powrót do miasta startowego
	totalCost += t.Matrix[currentCity][startCity]
	path = append(path, startCity)

	return Result{
		Path:     path,
		MinCost:  totalCost,
		Duration: time.Since(start),
	}
}