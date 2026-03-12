package main

import (
	"math"
	"math/rand"
	"time"
)

// SolveRandom implementuje algorytm losowy (Random Search).
// Przyjmuje jako parametr liczbę permutacji (iteracji) do wylosowania.
func (t TSPInstance) SolveRandom(permutations int) Result {
	start := time.Now()
	n := t.Size

	bestPath := make([]int, n+1)
	minCost := math.MaxInt32

	// Tworzymy tablicę miast do tasowania (bez miasta startowego 0)
	cities := make([]int, n-1)
	for i := 1; i < n; i++ {
		cities[i-1] = i
	}

	// Inicjujemy generator liczb pseudolosowych
	rand.Seed(time.Now().UnixNano())

	// Pętla wykonująca zadaną liczbę losowań
	for k := 0; k < permutations; k++ {
		// 1. Tasujemy tablicę miast
		rand.Shuffle(len(cities), func(i, j int) {
			cities[i], cities[j] = cities[j], cities[i]
		})

		// 2. Budujemy pełną trasę (od 0, przez przetasowane miasta, z powrotem do 0)
		currentPath := make([]int, n+1)
		currentPath[0] = 0
		copy(currentPath[1:n], cities)
		currentPath[n] = 0

		// 3. Obliczamy koszt tej konkretnej losowej trasy
		currentCost := 0
		for i := 0; i < n; i++ {
			from := currentPath[i]
			to := currentPath[i+1]
			currentCost += t.Matrix[from][to]
		}

		// 4. Jeśli wylosowana trasa jest najlepsza do tej pory, zapisujemy ją
		if currentCost < minCost {
			minCost = currentCost
			copy(bestPath, currentPath)
		}
	}

	return Result{
		Path:     bestPath,
		MinCost:  minCost,
		Duration: time.Since(start),
	}
}