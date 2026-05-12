package main

import (
	"math"
	"math/rand"
	"time"
)

type Result struct {
	Path     []int
	MinCost  int
	Duration time.Duration
}

type SimulatedAnnealing struct {
	Instance TSPInstance
	Config   SAConfig
}

func (sa *SimulatedAnnealing) Solve() Result {
	startTime := time.Now()
	maxDuration := time.Duration(sa.Config.MaxTimeMs) * time.Millisecond

	// 1. Inicjalizacja
	currentPath := sa.generateInitialSolution()
	currentCost := sa.Instance.CalculatePathCost(currentPath)

	bestPath := make([]int, sa.Instance.Size)
	copy(bestPath, currentPath)
	bestCost := currentCost

	// Jeden bufor wielokrotnie używany zamiast alokacji przy każdej iteracji
	neighborBuf := make([]int, sa.Instance.Size)

	temperature := sa.Config.InitialTemp

	// Pętla główna wyżarzania
	for {
		// Kryterium stopu - czas
		if time.Since(startTime) > maxDuration {
			break
		}

		// Kryterium stopu - niska temperatura 
		if temperature < 0.001 {
			break
		}

		for i := 0; i < sa.Config.EpochLength; i++ {
			// Generacja sąsiada i obliczenie kosztu
			_, neighborCost := sa.getNeighborIntoAndCost(currentPath, neighborBuf, currentCost)

			delta := neighborCost - currentCost

			if delta < 0 {
				copy(currentPath, neighborBuf)
				currentCost = neighborCost

				if currentCost < bestCost {
					copy(bestPath, currentPath)
					bestCost = currentCost
				}
			} else {
				probability := math.Exp(-float64(delta) / temperature)
				if rand.Float64() < probability {
					copy(currentPath, neighborBuf)
					currentCost = neighborCost
				}
			}
		}

		// Obniżanie temperatury według wybranego schematu
		temperature = sa.coolDown(temperature)
	}

	return Result{
		Path:     bestPath,
		MinCost:  bestCost,
		Duration: time.Since(startTime),
	}
}

func (sa *SimulatedAnnealing) coolDown(currentTemp float64) float64 {
	switch sa.Config.Cooling {
	case Geometric:
		return currentTemp * sa.Config.CoolingRate
	case Linear:
		newTemp := currentTemp - sa.Config.CoolingRate
		if newTemp < 0.0001 {
			return 0.0001
		}
		return newTemp
	case LundyMees:
		// T_{k+1} = T_k / (1 + beta * T_k)
		return currentTemp / (1.0 + sa.Config.CoolingRate*currentTemp)
	default:
		return currentTemp * 0.99
	}
}

func (sa *SimulatedAnnealing) generateInitialSolution() []int {
	path := make([]int, sa.Instance.Size)
	for i := 0; i < sa.Instance.Size; i++ {
		path[i] = i
	}

	if sa.Config.InitSol == RandomInit {
		rand.Shuffle(len(path), func(i, j int) {
			path[i], path[j] = path[j], path[i]
		})
	} else if sa.Config.InitSol == GreedyInit {
		// Zachłanny (Najbliższy Sąsiad) od losowego wierzchołka startowego
		startNode := rand.Intn(sa.Instance.Size)
		path[0] = startNode
		visited := make([]bool, sa.Instance.Size)
		visited[startNode] = true

		for i := 1; i < sa.Instance.Size; i++ {
			lastNode := path[i-1]
			bestNext := -1
			minDist := math.MaxInt32

			for j := 0; j < sa.Instance.Size; j++ {
				if lastNode == j {
					continue
				}
				if !visited[j] {
					dist := sa.Instance.Matrix[lastNode][j]
					if dist < minDist {
						minDist = dist
						bestNext = j
					}
				}
			}

			if bestNext != -1 {
				path[i] = bestNext
				visited[bestNext] = true
			} else {
				for j := 0; j < sa.Instance.Size; j++ {
					if !visited[j] {
						path[i] = j
						visited[j] = true
						break
					}
				}
			}
		}
	}

	return path
}

// getNeighborIntoAndCost generuje losowego sąsiada ścieżki currentPath i zapisuje wynik do bufora buf.
func (sa *SimulatedAnnealing) getNeighborIntoAndCost(currentPath []int, buf []int, currentCost int) ([]int, int) {
	copy(buf, currentPath)
	n := len(currentPath)
	i := rand.Intn(n)
	j := rand.Intn(n)

	for i == j {
		j = rand.Intn(n)
	}

	newCost := currentCost

	if sa.Config.NeighborGen == Swap {
		if i > j {
			i, j = j, i
		}
		buf[i], buf[j] = buf[j], buf[i]

		// O(1) Delta evaluation dla Swap
		if i == 0 || j == n-1 || j == i+1 {
			newCost = sa.Instance.CalculatePathCost(buf)
		} else {
			prevI, nextI := i-1, i+1
			prevJ, nextJ := j-1, j+1

			removed := sa.Instance.Matrix[currentPath[prevI]][currentPath[i]] +
				sa.Instance.Matrix[currentPath[i]][currentPath[nextI]] +
				sa.Instance.Matrix[currentPath[prevJ]][currentPath[j]] +
				sa.Instance.Matrix[currentPath[j]][currentPath[nextJ]]
			added := sa.Instance.Matrix[currentPath[prevI]][currentPath[j]] +
				sa.Instance.Matrix[currentPath[j]][currentPath[nextI]] +
				sa.Instance.Matrix[currentPath[prevJ]][currentPath[i]] +
				sa.Instance.Matrix[currentPath[i]][currentPath[nextJ]]
			newCost = currentCost - removed + added
		}

	} else if sa.Config.NeighborGen == Insert {
		elem := buf[i]
		if i < j {
			for k := i; k < j; k++ {
				buf[k] = buf[k+1]
			}
		} else {
			for k := i; k > j; k-- {
				buf[k] = buf[k-1]
			}
		}
		buf[j] = elem

		// O(1) Delta evaluation dla Insert
		if i == 0 || j == 0 || i == n-1 || j == n-1 || math.Abs(float64(i-j)) == 1 {
			newCost = sa.Instance.CalculatePathCost(buf)
		} else {
			prevI, nextI := i-1, i+1
			if i < j {
				// removed: prevI->i, i->nextI, j->j+1
				// added: prevI->nextI, j->i, i->j+1
				removed := sa.Instance.Matrix[currentPath[prevI]][currentPath[i]] +
					sa.Instance.Matrix[currentPath[i]][currentPath[nextI]] +
					sa.Instance.Matrix[currentPath[j]][currentPath[j+1]]
				added := sa.Instance.Matrix[currentPath[prevI]][currentPath[nextI]] +
					sa.Instance.Matrix[currentPath[j]][currentPath[i]] +
					sa.Instance.Matrix[currentPath[i]][currentPath[j+1]]
				newCost = currentCost - removed + added
			} else {
				// i > j
				// removed: prevI->i, i->nextI, j-1->j
				// added: prevI->nextI, j-1->i, i->j
				removed := sa.Instance.Matrix[currentPath[prevI]][currentPath[i]] +
					sa.Instance.Matrix[currentPath[i]][currentPath[nextI]] +
					sa.Instance.Matrix[currentPath[j-1]][currentPath[j]]
				added := sa.Instance.Matrix[currentPath[prevI]][currentPath[nextI]] +
					sa.Instance.Matrix[currentPath[j-1]][currentPath[i]] +
					sa.Instance.Matrix[currentPath[i]][currentPath[j]]
				newCost = currentCost - removed + added
			}
		}

	} 

	return buf, newCost
}

// Funkcja pomocnicza pozwalająca automatycznie dobrać temperaturę początkową
// na podstawie próbkowania przestrzeni rozwiązań i obliczenia średniej delty na plus
func (sa *SimulatedAnnealing) CalculateInitialTemp(prob float64, samples int) float64 {
	sumDelta := 0.0
	countPositiveDeltas := 0

	currentPath := sa.generateInitialSolution()
	currentCost := sa.Instance.CalculatePathCost(currentPath)
	// Jeden bufor dla próbkowania
	neighborBuf := make([]int, sa.Instance.Size)

	for k := 0; k < samples; k++ {
		neighbor, neighborCost := sa.getNeighborIntoAndCost(currentPath, neighborBuf, currentCost)
		delta := neighborCost - currentCost

		if delta > 0 {
			sumDelta += float64(delta)
			countPositiveDeltas++
		}

		copy(currentPath, neighbor)
		currentCost = neighborCost
	}

	if countPositiveDeltas == 0 {
		return 1000.0 // Default fallback
	}

	avgDelta := sumDelta / float64(countPositiveDeltas)

	// prob = e^(-avgDelta / T) => ln(prob) = -avgDelta / T => T = -avgDelta / ln(prob)
	initialTemp := -avgDelta / math.Log(prob)
	return initialTemp
}
