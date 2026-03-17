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

			// KROK 1: Koniec wycieczki - zamykamy pętlę i sprawdzamy wynik
			if len(path) == n {
				total := totalCost + t.Matrix[currentCity][startCity]
				
				if total < minTotalCost {
					minTotalCost = total
					
					// 1. Tworzymy pustą listę o rozmiarze trasy + 1 miejsce na miasto powrotne
					bestPath = make([]int, len(path)+1)
					// 2. Kopiujemy dotychczasowe miasta do nowej listy
					copy(bestPath, path)
					// 3. Wrzucamy miasto powrotne na sam koniec (na ostatni, pusty indeks)
					bestPath[len(path)] = startCity
				}
				return
			}

			// KROK 2: Rozglądamy się po mapie - szukamy najtańszego kosztu przejścia
			minEdgeCost := math.MaxInt32
			for j := 0; j < n; j++ {
				if !visited[j] && t.Matrix[currentCity][j] < minEdgeCost {
					minEdgeCost = t.Matrix[currentCity][j]
				}
			}

			// KROK 3: Sprawdzamy wszystkie odnogi z najtańszym kosztem
			for nextCity := 0; nextCity < n; nextCity++ {
				if !visited[nextCity] && t.Matrix[currentCity][nextCity] == minEdgeCost {
					visited[nextCity] = true // Zaznaczamy
					step(nextCity, append(path, nextCity), totalCost+minEdgeCost) // Idziemy głębiej
					visited[nextCity] = false // BACKTRACKING: Odznaczamy (cofamy decyzję)
				}
			}
		}

		// Startujemy z pierwszego wybranego miasta
		step(startCity, []int{startCity}, 0)
	}

	return Result{
		Path:     bestPath,
		MinCost:  minTotalCost,
		Duration: time.Since(start),
	}
}