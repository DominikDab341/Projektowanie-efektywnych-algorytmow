package main

import (
	"container/heap"
	"time"
)

const INF = int(^uint(0) >> 1)

// Węzeł w drzewie stanu 
type Node struct {
	Level   int
	Path    []int
	Cost    int
	Bound   int
	Visited uint64 // Bitmaska odwiedzonych miast (do 64 miast)
}

// Kolejka priorytetowa dla algorytmu Best-First-Search
type PriorityQueue []*Node

func (pq PriorityQueue) Len() int           { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool { return pq[i].Bound < pq[j].Bound }
func (pq PriorityQueue) Swap(i, j int)      { pq[i], pq[j] = pq[j], pq[i] }

func (pq *PriorityQueue) Push(x any) {
	item := x.(*Node)
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

// isVisited sprawdza czy miasto jest odwiedzone w bitmaskach
func isVisited(visited uint64, city int) bool {
	return visited&(1<<uint(city)) != 0
}

// setVisited ustawia miasto jako odwiedzone w bitmaskach
func setVisited(visited uint64, city int) uint64 {
	return visited | (1 << uint(city))
}

// calculateLowerBound wylicza dolne ograniczenie metodą redukcji macierzy kosztów
func calculateLowerBound(matrix [][]int, visited uint64, lastCity int, pathCost int, size int) int {
	
	// Zbieramy nieodwiedzone miasta
	unvisited := make([]int, 0, size)
	for i := 0; i < size; i++ {
		if !isVisited(visited, i) {
			unvisited = append(unvisited, i)
		}
	}

	k := len(unvisited)

	// Jeśli wszystkie miasta odwiedzone — zamykamy cykl
	if k == 0 {
		if matrix[lastCity][0] == -1 {
			return INF
		}
		return pathCost + matrix[lastCity][0]
	}

	// Budujemy podmacierz do redukcji:
	// Wiersze (skąd): lastCity + każde nieodwiedzone miasto
	// Kolumny (dokąd): każde nieodwiedzone miasto + miasto 0 (powrót)
	n := k + 1

	rowCities := make([]int, n)
	rowCities[0] = lastCity
	for i, c := range unvisited {
		rowCities[i+1] = c
	}

	colCities := make([]int, n)
	copy(colCities, unvisited)
	colCities[k] = 0

	// Wypełniamy podmacierz kosztami
	sub := make([][]int, n)
	for i := 0; i < n; i++ {
		sub[i] = make([]int, n)
		for j := 0; j < n; j++ {
			from := rowCities[i]
			to := colCities[j]
			if from == to || matrix[from][to] == -1 {
				sub[i][j] = INF
			} else {
				sub[i][j] = matrix[from][to]
			}
		}
	}

	reductionSum := 0

	// Redukcja WIERSZY — od każdego wiersza odejmujemy jego minimum
	for i := 0; i < n; i++ {
		minVal := INF
		for j := 0; j < n; j++ {
			if sub[i][j] < minVal {
				minVal = sub[i][j]
			}
		}
		if minVal == INF {
			return INF
		}
		if minVal > 0 {
			reductionSum += minVal
			for j := 0; j < n; j++ {
				if sub[i][j] != INF {
					sub[i][j] -= minVal
				}
			}
		}
	}

	// Redukcja KOLUMN — od każdej kolumny odejmujemy jej minimum
	for j := 0; j < n; j++ {
		minVal := INF
		for i := 0; i < n; i++ {
			if sub[i][j] < minVal {
				minVal = sub[i][j]
			}
		}
		if minVal == INF {
			return INF // Ślepy zaułek — kolumna bez żadnej krawędzi
		}
		if minVal > 0 {
			reductionSum += minVal
		}
	}

	return pathCost + reductionSum
}

// wyznaczenie początkowego rozwiązania heurystycznego (Nearest Neighbor)
func (t TSPInstance) getInitialBoundNN() (int, []int) {
	size := t.Size
	visited := make([]bool, size)
	path := []int{0}
	visited[0] = true
	cost := 0
	current := 0

	for step := 1; step < size; step++ {
		minEdge := INF
		next := -1
		for j := 0; j < size; j++ {
			if !visited[j] && t.Matrix[current][j] != -1 && t.Matrix[current][j] < minEdge {
				minEdge = t.Matrix[current][j]
				next = j
			}
		}
		if next == -1 {
			return INF, nil
		}
		cost += minEdge
		current = next
		path = append(path, current)
		visited[current] = true
	}

	if t.Matrix[current][0] == -1 {
		return INF, nil
	}
	cost += t.Matrix[current][0]
	path = append(path, 0)

	return cost, path
}

// Klonowanie plastra chroniące przed konfliktami w pamięci
func clonePath(p []int) []int {
	cp := make([]int, len(p))
	copy(cp, p)
	return cp
}


func (t TSPInstance) SolveBranchAndBound(metoda string, mode string, limitCzasu time.Duration) Result {
	resultChan := make(chan Result, 1)
	done := make(chan struct{})

	go func() {
		defer close(resultChan)
		start := time.Now() // Pomiar czasu wewnątrz goroutyny

		// Ustalenie Globalnego Najmniejszego Kosztu na START
		globalMinCost := INF
		var bestPath []int

		if mode == "NN" {
			c, p := t.getInitialBoundNN()
			if p != nil && c < globalMinCost {
				globalMinCost = c
				bestPath = p
			}
		}

		// Obliczenie Bound korzenia
		rootVisited := setVisited(0, 0)
		rootBound := calculateLowerBound(t.Matrix, rootVisited, 0, 0, t.Size)

		root := &Node{
			Level:   1,
			Path:    []int{0},
			Cost:    0,
			Bound:   rootBound,
			Visited: rootVisited,
		}

		timeoutTicks := 0

		if metoda == "BEST" {
			pq := make(PriorityQueue, 0)
			heap.Init(&pq)
			heap.Push(&pq, root)

			for pq.Len() > 0 {
				timeoutTicks++
				if timeoutTicks > 10000 {
					select {
					case <-done:
						return
					default:
						timeoutTicks = 0
					}
				}

				current := heap.Pop(&pq).(*Node)

				// Pruning węzła
				if current.Bound >= globalMinCost {
					continue
				}

				// Rozbudowa węzła
				lastCity := current.Path[len(current.Path)-1]

				if current.Level == t.Size {
					// Zamknięcie cyklu powrotem do startu
					if t.Matrix[lastCity][0] != -1 {
						finalCost := current.Cost + t.Matrix[lastCity][0]
						if finalCost < globalMinCost {
							globalMinCost = finalCost
							finalPath := clonePath(current.Path)
							finalPath = append(finalPath, 0)
							bestPath = finalPath
						}
					}
					continue
				}

				for i := 0; i < t.Size; i++ {
					if t.Matrix[lastCity][i] != -1 && !isVisited(current.Visited, i) {
						newCost := current.Cost + t.Matrix[lastCity][i]
						newVisited := setVisited(current.Visited, i)
						newBound := calculateLowerBound(t.Matrix, newVisited, i, newCost, t.Size)

						if newBound < globalMinCost {
							child := &Node{
								Level:   current.Level + 1,
								Path:    append(clonePath(current.Path), i),
								Cost:    newCost,
								Bound:   newBound,
								Visited: newVisited,
							}
							heap.Push(&pq, child)
						}
					}
				}
			}

		} else if metoda == "BREADTH" {
			// Klasyczna Kolejka FIFO
			queue := []*Node{root}

			for len(queue) > 0 {
				timeoutTicks++
				if timeoutTicks > 10000 {
					select {
					case <-done:
						return
					default:
						timeoutTicks = 0
					}
				}

				current := queue[0]
				queue = queue[1:]

				if current.Bound >= globalMinCost {
					continue
				}

				lastCity := current.Path[len(current.Path)-1]

				if current.Level == t.Size {
					if t.Matrix[lastCity][0] != -1 {
						finalCost := current.Cost + t.Matrix[lastCity][0]
						if finalCost < globalMinCost {
							globalMinCost = finalCost
							finalPath := clonePath(current.Path)
							finalPath = append(finalPath, 0)
							bestPath = finalPath
						}
					}
					continue
				}

				for i := 0; i < t.Size; i++ {
					if t.Matrix[lastCity][i] != -1 && !isVisited(current.Visited, i) {
						newCost := current.Cost + t.Matrix[lastCity][i]
						newVisited := setVisited(current.Visited, i)
						newBound := calculateLowerBound(t.Matrix, newVisited, i, newCost, t.Size)

						if newBound < globalMinCost {
							child := &Node{
								Level:   current.Level + 1,
								Path:    append(clonePath(current.Path), i),
								Cost:    newCost,
								Bound:   newBound,
								Visited: newVisited,
							}
							queue = append(queue, child)
						}
					}
				}
			}
		}

		resultChan <- Result{
			Path:     bestPath,
			MinCost:  globalMinCost,
			Duration: time.Since(start),
		}
	}()

	if limitCzasu > 0 {
		select {
		case res := <-resultChan:
			return res
		case <-time.After(limitCzasu):
			close(done)
			return Result{MinCost: -1, Duration: limitCzasu}
		}
	} else {
		res := <-resultChan
		return res
	}
}
