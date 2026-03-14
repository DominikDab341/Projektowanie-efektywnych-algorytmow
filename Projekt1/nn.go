package main

import (
	"math"
	"time"
)

// SolveNN implementuje algorytm Nearest Neighbor z obsługą rozgałęzień (remisów)
func (t TSPInstance) SolveNN(startCity int) Result {
	start := time.Now()
	n := t.Size

	var bestPath []int
	minTotalCost := math.MaxInt32

	// Wewnętrzna funkcja rekurencyjna do obsługi rozgałęzień
	var dfs func(currentCity int, visited []bool, currentPath []int, currentCost, step int)
	dfs = func(currentCity int, visited []bool, currentPath []int, currentCost, step int) {
		
		// Warunek końcowy: odwiedzono wszystkie miasta
		if step == n {
			// Dodajemy powrót do bazy
			total := currentCost + t.Matrix[currentCity][startCity]
			if total < minTotalCost {
				minTotalCost = total
				
				// Zapisujemy nową najlepszą trasę
				bestPath = make([]int, len(currentPath)+1)
				copy(bestPath, currentPath)
				bestPath[len(currentPath)] = startCity
			}
			return
		}

		minEdgeCost := math.MaxInt32
		var nextCities []int

		// Szukamy najmniejszego kosztu, ale zapamiętujemy WSZYSTKIE miasta, które go mają
		for j := 0; j < n; j++ {
			if !visited[j] {
				cost := t.Matrix[currentCity][j]
				if cost < minEdgeCost {
					// Znaleziono nowy, mniejszy koszt - resetujemy listę remisów
					minEdgeCost = cost
					nextCities = []int{j}
				} else if cost == minEdgeCost {
					// Znaleziono miasto o identycznym, najmniejszym koszcie - dodajemy do listy
					nextCities = append(nextCities, j)
				}
			}
		}

		// Rozgałęzienie: idziemy do KAŻDEGO miasta z listy najtańszych
		for _, nextCity := range nextCities {
			visited[nextCity] = true
			currentPath = append(currentPath, nextCity)
			
			// Wywołanie rekurencyjne dla wybranej odnogi
			dfs(nextCity, visited, currentPath, currentCost+minEdgeCost, step+1)
			
			// Backtracking: cofamy się, aby sprawdzić inne opcje z listy remisów
			visited[nextCity] = false
			currentPath = currentPath[:len(currentPath)-1]
		}
	}

	// Inicjalizacja stanu początkowego i start rekurencji
	visited := make([]bool, n)
	visited[startCity] = true
	dfs(startCity, visited, []int{startCity}, 0, 1)

	return Result{
		Path:     bestPath,
		MinCost:  minTotalCost,
		Duration: time.Since(start),
	}
}