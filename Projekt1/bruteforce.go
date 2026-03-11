package main

import (
	"math"
	"time"
)

// SolveBruteForce implementuje algorytm przeglądu zupełnego przy użyciu rekurencji
func (t TSPInstance) SolveBruteForce() Result {
	start := time.Now()
	n := t.Size

	visited := make([]bool, n)
	currentPath := make([]int, n+1)
	bestPath := make([]int, n+1)
	minCost := math.MaxInt32

	// Zawsze zaczynamy od miasta 0
	visited[0] = true
	currentPath[0] = 0

	// Rekurencyjna funkcja wewnętrzna (backtracking)
	var backtrack func(currCity, count, currentCost int)
	backtrack = func(currCity, count, currentCost int) {
		// Warunek końcowy: odwiedziliśmy wszystkie N miast
		if count == n {
			// Dodajemy koszt powrotu do miasta startowego (0)
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
				// Krok w przód (wybór miasta)
				visited[nextCity] = true
				currentPath[count] = nextCity

				// Wywołanie rekurencyjne dla kolejnego kroku
				backtrack(nextCity, count+1, currentCost+t.Matrix[currCity][nextCity])

				// Cofnięcie zmian (backtracking) - powrót, aby sprawdzić inne ścieżki
				visited[nextCity] = false
			}
		}
	}

	// Uruchomienie rekurencji: startujemy w mieście 0, licznik odwiedzonych to 1, koszt początkowy 0
	backtrack(0, 1, 0)

	return Result{
		Path:     bestPath,
		MinCost:  minCost,
		Duration: time.Since(start),
	}
}