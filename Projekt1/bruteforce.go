package main

import (
	"math"
	"time"
)

func (t TSPInstance) SolveBruteForce() Result {
	start := time.Now()
	n := t.Size

	visited := make([]bool, n)
	currentPath := make([]int, n+1)
	bestPath := make([]int, n+1)
	minCost := math.MaxInt32


	visited[0] = true
	currentPath[0] = 0

	// Rekurencyjna funkcja wewnętrzna 
	var backtrack func(currCity, count, currentCost int)
	backtrack = func(currCity, count, currentCost int) {
		// Warunek końcowy: odwiedziliśmy wszystkie N miast
		if count == n {
			totalCost := currentCost + t.Matrix[currCity][0]
			
			// Aktualizacja najlepszego wyniku
			if totalCost < minCost {
				minCost = totalCost
				copy(bestPath, currentPath) 
				bestPath[n] = 0             
			}
			return
		}

		// Przeszukujemy wszystkie możliwe kolejne miasta
		for nextCity := 0; nextCity < n; nextCity++ {
			if !visited[nextCity] {
				// Krok w przód 
				visited[nextCity] = true
				currentPath[count] = nextCity

				// Wywołanie rekurencyjne dla kolejnego kroku
				backtrack(nextCity, count+1, currentCost+t.Matrix[currCity][nextCity])

				// Powrót, aby sprawdzić inne ścieżki
				visited[nextCity] = false
			}
		}
	}


	backtrack(0, 1, 0)

	return Result{
		Path:     bestPath,
		MinCost:  minCost,
		Duration: time.Since(start),
	}
}