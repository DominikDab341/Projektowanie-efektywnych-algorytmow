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

		// Ta funkcja zastępuje "for step" z NN.
		// Zamiast jednej pętli, wywołuje samą siebie (rekurencja),
		// dzięki czemu może się cofnąć i spróbować innej drogi (backtracking).
		var step func(currentCity int, path []int, totalCost int)
		step = func(currentCity int, path []int, totalCost int) {

			// Wszystkie miasta odwiedzone -> powrót do bazy (tak jak po pętli w NN)
			if len(path) == n {
				total := totalCost + t.Matrix[currentCity][startCity]
				if total < minTotalCost {
					minTotalCost = total
					bestPath = append([]int{}, append(path, startCity)...)
				}
				return
			}

			// === IDENTYCZNIE JAK W NN: szukamy minEdgeCost ===
			minEdgeCost := math.MaxInt32
			for j := 0; j < n; j++ {
				if !visited[j] {
					cost := t.Matrix[currentCity][j]
					if cost < minEdgeCost {
						minEdgeCost = cost
					}
				}
			}

			// === PRAWIE JAK W NN: idziemy do nextCity ===
			// Różnica: w NN bierzemy JEDNO miasto, tu bierzemy KAŻDE z kosztem == minEdgeCost
			for nextCity := 0; nextCity < n; nextCity++ {
				if !visited[nextCity] && t.Matrix[currentCity][nextCity] == minEdgeCost {
					visited[nextCity] = true
					step(nextCity, append(path, nextCity), totalCost+minEdgeCost) // zamiast "for step++"
					visited[nextCity] = false // BACKTRACKING - cofamy decyzję
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