package main

import (
	"math"
	"time"
)

// SolveNN implementuje algorytm Najbliższego Sąsiada (Nearest Neighbor)
// Przyjmuje miasto startowe jako parametr, co ułatwi późniejszą implementację RNN.
func (t TSPInstance) SolveNN(startCity int) Result {
	start := time.Now()
	n := t.Size

	visited := make([]bool, n)
	path := make([]int, 0, n+1) // Pojemność n+1, bo wracamy do startu
	totalCost := 0

	// Ustawiamy się w mieście początkowym
	currentCity := startCity
	visited[currentCity] = true
	path = append(path, currentCity)

	// Musimy odwiedzić jeszcze n-1 miast
	for step := 1; step < n; step++ {
		nextCity := -1
		minEdgeCost := math.MaxInt32

		// Szukamy najtańszego połączenia do nieodwiedzonego miasta
		for j := 0; j < n; j++ {
			if !visited[j] {
				cost := t.Matrix[currentCity][j]
				if cost < minEdgeCost {
					minEdgeCost = cost
					nextCity = j
				}
			}
		}

		// Przechodzimy do znalezionego miasta
		visited[nextCity] = true
		path = append(path, nextCity)
		totalCost += minEdgeCost
		currentCity = nextCity
	}

	// Po odwiedzeniu wszystkich miast, wracamy do miasta startowego
	totalCost += t.Matrix[currentCity][startCity]
	path = append(path, startCity)

	return Result{
		Path:     path,
		MinCost:  totalCost,
		Duration: time.Since(start),
	}
}