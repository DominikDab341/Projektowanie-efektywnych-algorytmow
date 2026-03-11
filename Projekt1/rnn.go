package main

import (
	"math"
	"time"
)

// SolveRNN implementuje algorytm Wielokrotnego Najbliższego Sąsiada (Repetitive Nearest Neighbor)
func (t TSPInstance) SolveRNN() Result {
	start := time.Now()
	n := t.Size

	var bestPath []int
	minCost := math.MaxInt32

	// Uruchamiamy NN dla każdego możliwego miasta startowego (od 0 do N-1)
	for startCity := 0; startCity < n; startCity++ {
		
		// Wywołujemy algorytm NN dla aktualnego miasta startowego
		currentResult := t.SolveNN(startCity)

		// Jeśli znaleziona trasa jest lepsza (krótsza) niż dotychczasowa, zapamiętujemy ją
		if currentResult.MinCost < minCost {
			minCost = currentResult.MinCost
			
			// Tworzymy nową tablicę i kopiujemy do niej ścieżkę,
			// aby uniknąć nadpisania jej w kolejnych iteracjach pętli
			bestPath = make([]int, len(currentResult.Path))
			copy(bestPath, currentResult.Path)
		}
	}

	return Result{
		Path:     bestPath,
		MinCost:  minCost,
		// Czas mierzymy od początku funkcji, obejmując wszystkie N wywołań NN
		Duration: time.Since(start),
	}
}