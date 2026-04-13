package main

import (
	"fmt"
	"runtime"
	"time"
)

// RunAutomatedTests przeprowadza statystyczną analizę działania Branch&Bound ATSP wg wymogów zadania 2.
// Generuje instancje RAZ, a potem testuje każdy wariant na TEJ SAMEJ puli — uczciwe porównanie.
func RunAutomatedTests() {
	LimitCzasu := 5 * time.Minute
	LimitInstancji := 100

	fmt.Println("\n========================================================================")
	fmt.Println("       BADANIA BRANCH & BOUND — ATSP (Limit: 5 min / instancję)")
	fmt.Println("========================================================================")

	// =====================================================================
	// TEST 1: BREADTH-FIRST-SEARCH (INF) — N od 8 do 12
	// =====================================================================
	// sizesBFS := []int{8, 9, 10, 11, 12}

	// fmt.Println("\n--- [1/3] BREADTH-FIRST-SEARCH (Ograniczenie: INF) ---")
	// fmt.Println("N\tŚredni czas\t\tTimeouty")

	// for _, n := range sizesBFS {
	// 	// Generujemy instancje RAZ
	// 	instances := make([]TSPInstance, LimitInstancji)
	// 	for i := 0; i < LimitInstancji; i++ {
	// 		instances[i] = GenerateRandom(n)
	// 	}

	// 	var sumTime time.Duration
	// 	timeouts := 0

	// 	for i := 0; i < LimitInstancji; i++ {
	// 		res := instances[i].SolveBranchAndBound("BREADTH", "INF", LimitCzasu)
	// 		if res.MinCost == -1 {
	// 			timeouts++
	// 		} else {
	// 			sumTime += res.Duration
	// 		}
	// 	}

	// 	sukcesy := LimitInstancji - timeouts
	// 	var avg time.Duration
	// 	if sukcesy > 0 {
	// 		avg = time.Duration(float64(sumTime.Nanoseconds()) / float64(sukcesy))
	// 	} else {
	// 		avg = LimitCzasu
	// 	}

	// 	fmt.Printf("%d\t%v\t\t%d%%\n", n, avg, timeouts)

	// 	if timeouts == LimitInstancji {
	// 		fmt.Printf(">> BFS(INF): Blokada przy N=%d. Wszystkie instancje przekroczyły limit.\n", n)
	// 		break
	// 	}
	// }

	// =====================================================================
	// TEST 2 & 3: BEST-FIRST-SEARCH (INF vs NN) — te same instancje!
	// =====================================================================
	sizesBest := []int{8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21 ,22 ,23 ,24, 25, 26, 27, 28}

	fmt.Println("\n--- [2/3 + 3/3] BEST-FIRST-SEARCH (INF vs NN — te same instancje) ---")
	fmt.Println("N\tBest(INF) czas\t\tINF T/O\t\tBest(NN) czas\t\tNN T/O")

	for _, n := range sizesBest {
		// Generujemy instancje RAZ — obie metody dostaną DOKŁADNIE TE SAME dane
		instances := make([]TSPInstance, LimitInstancji)
		for i := 0; i < LimitInstancji; i++ {
			instances[i] = GenerateRandom(n)
		}

		var sumTimeINF, sumTimeNN time.Duration
		timeoutINF, timeoutNN := 0, 0

		runtime.GC() // Wymuszenie GC przed serią pomiarów — żeby nie strzelał w trakcie

		for i := 0; i < LimitInstancji; i++ {
			// Best(INF) na tej samej instancji
			resINF := instances[i].SolveBranchAndBound("BEST", "INF", LimitCzasu)
			if resINF.MinCost == -1 {
				timeoutINF++
			} else {
				sumTimeINF += resINF.Duration
			}

			// Best(NN) na tej samej instancji
			resNN := instances[i].SolveBranchAndBound("BEST", "NN", LimitCzasu)
			if resNN.MinCost == -1 {
				timeoutNN++
			} else {
				sumTimeNN += resNN.Duration
			}
		}

		sukcesyINF := LimitInstancji - timeoutINF
		sukcesyNN := LimitInstancji - timeoutNN

		var avgINF, avgNN time.Duration
		if sukcesyINF > 0 {
			avgINF = time.Duration(float64(sumTimeINF.Nanoseconds()) / float64(sukcesyINF))
		} else {
			avgINF = LimitCzasu
		}
		if sukcesyNN > 0 {
			avgNN = time.Duration(float64(sumTimeNN.Nanoseconds()) / float64(sukcesyNN))
		} else {
			avgNN = LimitCzasu
		}

		fmt.Printf("%d\t%v\t\t%d%%\t\t%v\t\t%d%%\n",
			n, avgINF, timeoutINF, avgNN, timeoutNN)

		if timeoutINF == LimitInstancji && timeoutNN == LimitInstancji {
			fmt.Printf(">> Best: Blokada przy N=%d. Wszystkie instancje przekroczyły limit.\n", n)
			break
		}
	}

	fmt.Println("\n========================================================================")
	fmt.Println("                        BADANIA ZAKOŃCZONE")
	fmt.Println("========================================================================")
}
