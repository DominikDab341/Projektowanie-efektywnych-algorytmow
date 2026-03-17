package main

import (
	"math"
	"time"
)

func (t TSPInstance) SolveRNN() Result {
	start := time.Now()
	n := t.Size

	var bestPath []int
	minTotalCost := math.MaxInt32

	for startCity := 0; startCity < n; startCity++ {

		visited := make([]bool, n)
		visited[startCity] = true

		var step func(currentCity int, path []int, totalCost int)
		step = func(currentCity int, path []int, totalCost int) {

			// zamykamy pętlę i sprawdzamy wynik
			if len(path) == n {
				total := totalCost + t.Matrix[currentCity][startCity]
				
				if total < minTotalCost {
					minTotalCost = total
					
					bestPath = make([]int, len(path)+1)
					copy(bestPath, path)
					bestPath[len(path)] = startCity
				}
				return
			}

			// szukamy najtańszego kosztu przejścia
			minEdgeCost := math.MaxInt32
			for j := 0; j < n; j++ {
				if !visited[j] && t.Matrix[currentCity][j] < minEdgeCost {
					minEdgeCost = t.Matrix[currentCity][j]
				}
			}

			// Sprawdzamy wszystkie odnogi z najtańszym kosztem
			for nextCity := 0; nextCity < n; nextCity++ {
				if !visited[nextCity] && t.Matrix[currentCity][nextCity] == minEdgeCost {
					visited[nextCity] = true
					step(nextCity, append(path, nextCity), totalCost+minEdgeCost)
					visited[nextCity] = false
				}
			}
		}

		step(startCity, []int{startCity}, 0)
	}

	return Result{
		Path:     bestPath,
		MinCost:  minTotalCost,
		Duration: time.Since(start),
	}
}