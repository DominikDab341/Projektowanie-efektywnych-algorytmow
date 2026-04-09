package main

import (
	"container/heap"
	"time"
)

// Węzeł w drzewie stanu dla Branch & Bound
type Node struct {
	Level int
	Path  []int
	Cost  int
	Bound int
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

// calculateLowerBound wylicza dolne ograniczenie dla węzła (The Lower Bound)
func calculateLowerBound(matrix [][]int, path []int, size int) int {
	cost := 0
	visited := make([]bool, size)

	// Sumujemy koszty z aktualnej, już wybranej ścieżki
	for i := 0; i < len(path)-1; i++ {
		cost += matrix[path[i]][path[i+1]]
		visited[path[i]] = true
	}
	last := path[len(path)-1]
	visited[last] = true

	// Jeżeli graf jest już w pełni odwiedzony, wystarczy zwrócić koszt dotychczasowy z powrotem by zamknąć graf
	if len(path) == size {
		if matrix[last][path[0]] == -1 {
			return int(^uint(0) >> 1)
		}
		return cost + matrix[last][path[0]]
	}

	// Dodajemy najkrótsze możliwe ścieżki wychodzące z ostatnich i nieodwiedzonych wierzchołków.
	// Dla Ostatniego węzła z wybranej ścieżki:
	minOutLast := int(^uint(0) >> 1)
	for j := 0; j < size; j++ {
		if !visited[j] && matrix[last][j] != -1 && matrix[last][j] < minOutLast {
			minOutLast = matrix[last][j]
		}
	}
	if minOutLast == int(^uint(0)>>1) {
		return int(^uint(0) >> 1) // Ślepy zaułek
	}
	cost += minOutLast

	// Dla pozostałych wierzchołków NIEodwiedzonych z grafu:
	for i := 0; i < size; i++ {
		if !visited[i] {
			minOut := int(^uint(0) >> 1)
			for j := 0; j < size; j++ {
				// Możemy z nieodwiedzonego wyjść do następnego nieodwiedzonego,
				// LUB zamknąć cykl, zwracając się do wierzchołka 0.
				if i != j && matrix[i][j] != -1 && (!visited[j] || j == path[0]) {
					if matrix[i][j] < minOut {
						minOut = matrix[i][j]
					}
				}
			}
			if minOut == int(^uint(0)>>1) {
				return int(^uint(0) >> 1) // Ślepy zaułek grafu
			}
			cost += minOut
		}
	}

	return cost
}

// Funkcja pomocnicza: wyznaczenie początkowego rozwiązania heurystycznego (Nearest Neighbor)
// By służyło jako wstępnie narzucona bariera ucinania gałęzi (Pruning limit)
func (t TSPInstance) getInitialBoundNN() (int, []int) {
	size := t.Size
	visited := make([]bool, size)
	path := []int{0}
	visited[0] = true
	cost := 0
	current := 0

	for step := 1; step < size; step++ {
		minEdge := int(^uint(0) >> 1)
		next := -1
		for j := 0; j < size; j++ {
			if !visited[j] && t.Matrix[current][j] != -1 && t.Matrix[current][j] < minEdge {
				minEdge = t.Matrix[current][j]
				next = j
			}
		}
		if next == -1 {
			return int(^uint(0) >> 1), nil // Gdyby NN wylosował ślepok, oddaje nieskończoność
		}
		cost += minEdge
		current = next
		path = append(path, current)
		visited[current] = true
	}

	if t.Matrix[current][0] == -1 {
		return int(^uint(0) >> 1), nil
	}
	cost += t.Matrix[current][0]
	path = append(path, 0)

	return cost, path
}

// Klonowanie plastra chroniące przed konfliktami w pamięci współbieżnej
func clonePath(p []int) []int {
	cp := make([]int, len(p))
	copy(cp, p)
	return cp
}

func contains(slice []int, item int) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}

// SolveBranchAndBound rozwiązuje instancję ATSP bazując na metodzie B&B ("BREADTH", "BEST").
// 'mode' definiuje użyty limiter początkowy: "NN" (Nearest Neighbor) lub "INF" (Infinity).
func (t TSPInstance) SolveBranchAndBound(metoda string, mode string, limitCzasu time.Duration) Result {
	start := time.Now()
	resultChan := make(chan Result, 1) // Buforowana by gorutyna się wyłączyła
	done := make(chan struct{})

	go func() {
		defer close(resultChan)

		// Ustalenie Globalnego Najmniejszego Kosztu na START (Upper Bound Limit)
		globalMinCost := int(^uint(0) >> 1)
		var bestPath []int

		if mode == "NN" {
			c, p := t.getInitialBoundNN()
			if p != nil && c < globalMinCost {
				globalMinCost = c
				bestPath = p
			}
		}

		// Obliczenie Bound korzenia
		rootPath := []int{0}
		rootBound := calculateLowerBound(t.Matrix, rootPath, t.Size)

		root := &Node{
			Level: 1,
			Path:  rootPath,
			Cost:  0,
			Bound: rootBound,
		}

		// Zabezpieczenie przed timeoutami podczas kręcenia głębokiej pętli
		timeoutTicks := 0

		if metoda == "BEST" {
			pq := make(PriorityQueue, 0)
			heap.Init(&pq)
			heap.Push(&pq, root)

			for pq.Len() > 0 {
				// Bezpiecznik pozwalający po 10 tys węzłów rzucić sprawdzenie czy główny thread
				// wysłał na zewnątrz sygnał TimeOut by zamknąć gorutynę obliczeniową.
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

				// Pruning węzła (jeżeli teoretycznie możliwe najlepiej wyliczone
				// z niego ścieżki i tak są wyższe niż nasz obecny absolutny rekord to Ucinamy Gałąź)
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
							var finalPath = clonePath(current.Path)
							finalPath = append(finalPath, 0)
							bestPath = finalPath
						}
					}
					continue
				}

				for i := 0; i < t.Size; i++ {
					if t.Matrix[lastCity][i] != -1 && !contains(current.Path, i) {
						newPath := append(clonePath(current.Path), i)
						newCost := current.Cost + t.Matrix[lastCity][i]
						newBound := calculateLowerBound(t.Matrix, newPath, t.Size)

						if newBound < globalMinCost {
							child := &Node{
								Level: current.Level + 1,
								Path:  newPath,
								Cost:  newCost,
								Bound: newBound,
							}
							heap.Push(&pq, child)
						}
					}
				}
			}

		} else if metoda == "BREADTH" {
			// Klasyczna Kolejka FIFO (Plaster)
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
							var finalPath = clonePath(current.Path)
							finalPath = append(finalPath, 0)
							bestPath = finalPath
						}
					}
					continue
				}

				for i := 0; i < t.Size; i++ {
					if t.Matrix[lastCity][i] != -1 && !contains(current.Path, i) {
						newPath := append(clonePath(current.Path), i)
						newCost := current.Cost + t.Matrix[lastCity][i]
						newBound := calculateLowerBound(t.Matrix, newPath, t.Size)

						// Ograniczanie też dodajemy by oszczędzić trochę w BFS, ale samo wrzucanie
						// będzie bez priorytetów (od nogi, ślepe przeglądanie poziom za poziomem)
						if newBound < globalMinCost {
							child := &Node{
								Level: current.Level + 1,
								Path:  newPath,
								Cost:  newCost,
								Bound: newBound,
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
			// Timeout algorytmu z powodu braku w pamięci/konsekwencji
			close(done) // Sygnał do zabicia gorutyny by nie kradła RAMu i CPU.
			return Result{MinCost: -1, Duration: limitCzasu}
		}
	} else {
		// Tryb Menu (bez timeoutu)
		res := <-resultChan
		return res
	}
}
