package main

import (
	"math"
	"time"
)

// SolveRNN - Repetitive Nearest Neighbor z rozgałęzieniami (obsługa remisów)
func (t TSPInstance) SolveRNN() Result {
	start := time.Now()
	n := t.Size

	var bestPath []int
	minCost := math.MaxInt32

	for startCity := 0; startCity < n; startCity++ {
		visited := make([]bool, n)
		visited[startCity] = true

		var dfs func(city int, path []int, cost, step int)
		dfs = func(city int, path []int, cost, step int) {
			if step == n {
				total := cost + t.Matrix[city][startCity]
				if total < minCost {
					minCost = total
					bestPath = append([]int{}, append(path, startCity)...)
				}
				return
			}

			// Znajdź minimalny koszt krawędzi
			minEdge := math.MaxInt32
			for j := 0; j < n; j++ {
				if !visited[j] && t.Matrix[city][j] < minEdge {
					minEdge = t.Matrix[city][j]
				}
			}

			// Rozgałęzienie: idź do każdego miasta z tym minimalnym kosztem
			for j := 0; j < n; j++ {
				if !visited[j] && t.Matrix[city][j] == minEdge {
					visited[j] = true
					dfs(j, append(path, j), cost+minEdge, step+1)
					visited[j] = false
				}
			}
		}

		dfs(startCity, []int{startCity}, 0, 1)
	}

	return Result{
		Path:     bestPath,
		MinCost:  minCost,
		Duration: time.Since(start),
	}
}