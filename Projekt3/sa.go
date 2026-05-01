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

	temperature := sa.Config.InitialTemp
	epochs := 0

	// Pętla główna wyżarzania
	for {
		// Kryterium stopu - czas
		if time.Since(startTime) > maxDuration {
			break
		}

		// Kryterium stopu - niska temperatura (np. poniżej 0.001)
		if temperature < 0.001 {
			break
		}

		for i := 0; i < sa.Config.EpochLength; i++ {
			// Generacja sąsiada
			neighbor := sa.getNeighbor(currentPath)
			neighborCost := sa.Instance.CalculatePathCost(neighbor)

			// Akceptacja lub odrzucenie
			delta := neighborCost - currentCost

			if delta < 0 {
				// Zawsze akceptujemy lepsze rozwiązanie
				currentPath = neighbor
				currentCost = neighborCost

				if currentCost < bestCost {
					copy(bestPath, currentPath)
					bestCost = currentCost
				}
			} else {
				// Akceptacja gorszego rozwiązania z prawdopodobieństwem e^(-delta/T)
				probability := math.Exp(-float64(delta) / temperature)
				if rand.Float64() < probability {
					currentPath = neighbor
					currentCost = neighborCost
				}
			}
		}

		// Obniżanie temperatury według schematu
		epochs++
		temperature = sa.coolDown(temperature, epochs)
	}

	return Result{
		Path:     bestPath,
		MinCost:  bestCost,
		Duration: time.Since(startTime),
	}
}

func (sa *SimulatedAnnealing) coolDown(currentTemp float64, epoch int) float64 {
	switch sa.Config.Cooling {
	case Geometric:
		return currentTemp * sa.Config.CoolingRate
	case Linear:
		// Linear może szybko spaść poniżej 0, więc ograniczamy
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
				if !visited[j] && lastNode != j {
					// W niektórych ATSP dystans może być ujemny na przekątnej (np -1), ignorujemy
					dist := sa.Instance.Matrix[lastNode][j]
					if dist < minDist && dist >= 0 {
						minDist = dist
						bestNext = j
					}
				}
			}

			if bestNext != -1 {
				path[i] = bestNext
				visited[bestNext] = true
			} else {
				// Jeśli nie udało się znaleźć (bardzo dziwny graf), to bierzemy cokolwiek
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

func (sa *SimulatedAnnealing) getNeighbor(currentPath []int) []int {
	neighbor := make([]int, len(currentPath))
	copy(neighbor, currentPath)

	n := len(currentPath)
	i := rand.Intn(n)
	j := rand.Intn(n)

	for i == j {
		j = rand.Intn(n)
	}

	if sa.Config.NeighborGen == Swap {
		neighbor[i], neighbor[j] = neighbor[j], neighbor[i]
	} else if sa.Config.NeighborGen == Invert {
		if i > j {
			i, j = j, i
		}
		// Odwrócenie podciągu od i do j włącznie
		for k := 0; k < (j-i+1)/2; k++ {
			neighbor[i+k], neighbor[j-k] = neighbor[j-k], neighbor[i+k]
		}
	}

	return neighbor
}

// Funkcja pomocnicza pozwalająca automatycznie dobrać temperaturę początkową
// na podstawie próbkowania przestrzeni rozwiązań i obliczenia średniej delty na plus
func (sa *SimulatedAnnealing) CalculateInitialTemp(prob float64, samples int) float64 {
	sumDelta := 0.0
	countPositiveDeltas := 0

	currentPath := sa.generateInitialSolution()
	currentCost := sa.Instance.CalculatePathCost(currentPath)

	for k := 0; k < samples; k++ {
		neighbor := sa.getNeighbor(currentPath)
		neighborCost := sa.Instance.CalculatePathCost(neighbor)
		delta := neighborCost - currentCost

		if delta > 0 {
			sumDelta += float64(delta)
			countPositiveDeltas++
		}

		currentPath = neighbor
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
