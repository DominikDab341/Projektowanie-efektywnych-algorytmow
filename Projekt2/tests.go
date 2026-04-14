package main

import (
	"fmt"
	"runtime"
	"time"
)

// RunAutomatedTests przeprowadza testy algorytmu
func RunAutomatedTests() {
	LimitCzasu := 2 * time.Minute
	LimitInstancji := 100

	fmt.Println("\n========================================================================")
	fmt.Println("       BADANIA BRANCH & BOUND — ATSP (Limit: 2 min / instancję)")
	fmt.Println("========================================================================")

	sizesBFS := []int{8, 9, 10, 11}

	fmt.Println("\n--- [1/3] BREADTH-FIRST-SEARCH (Ograniczenie: INF) ---")
	fmt.Println("N\tŚredni czas\t\tTimeouty")

	for _, n := range sizesBFS {
		// Generujemy instancje RAZ
		instances := make([]TSPInstance, LimitInstancji)
		for i := 0; i < LimitInstancji; i++ {
			instances[i] = GenerateRandom(n)
		}

		var sumTime time.Duration
		timeouts := 0

		for i := 0; i < LimitInstancji; i++ {
			res := instances[i].SolveBranchAndBound("BREADTH", "INF", LimitCzasu)
			if res.MinCost == -1 {
				timeouts++
			} else {
				sumTime += res.Duration
			}
		}

		sukcesy := LimitInstancji - timeouts
		var avg time.Duration
		if sukcesy > 0 {
			avg = time.Duration(float64(sumTime.Nanoseconds()) / float64(sukcesy))
		} else {
			avg = LimitCzasu
		}

		fmt.Printf("%d\t%v\t\t%d%%\n", n, avg, timeouts)

		if timeouts == LimitInstancji {
			fmt.Printf(">> BFS(INF): Blokada przy N=%d. Wszystkie instancje przekroczyły limit.\n", n)
			break
		}
	}

	sizesBest := []int{8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40}

	fmt.Println("\n--- [2/3 + 3/3] BEST-FIRST-SEARCH (INF vs NN) ---")
	fmt.Println("N\tINF śr.(udane)\t\tINF śr.(wszystkie)\tINF T/O\t\tNN śr.(udane)\t\tNN śr.(wszystkie)\tNN T/O")

	for _, n := range sizesBest {
		instances := make([]TSPInstance, LimitInstancji)
		for i := 0; i < LimitInstancji; i++ {
			instances[i] = GenerateRandom(n)
		}

		var sumTimeINF, sumTimeNN time.Duration
		timeoutINF, timeoutNN := 0, 0

		runtime.GC()

		for i := 0; i < LimitInstancji; i++ {
			resINF := instances[i].SolveBranchAndBound("BEST", "INF", LimitCzasu)
			if resINF.MinCost == -1 {
				timeoutINF++
			} else {
				sumTimeINF += resINF.Duration
			}

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

		sumAllINF := sumTimeINF + time.Duration(timeoutINF)*LimitCzasu
		sumAllNN := sumTimeNN + time.Duration(timeoutNN)*LimitCzasu
		avgAllINF := time.Duration(float64(sumAllINF.Nanoseconds()) / float64(LimitInstancji))
		avgAllNN := time.Duration(float64(sumAllNN.Nanoseconds()) / float64(LimitInstancji))

		fmt.Printf("%d\t%v\t\t%v\t\t%d%%\t\t%v\t\t%v\t\t%d%%\n",
			n, avgINF, avgAllINF, timeoutINF*100/LimitInstancji, avgNN, avgAllNN, timeoutNN*100/LimitInstancji)

		if timeoutINF == LimitInstancji && timeoutNN == LimitInstancji {
			fmt.Printf(">> Best: Blokada przy N=%d. Wszystkie instancje przekroczyły limit.\n", n)
			break
		}
	}

	fmt.Println("\n========================================================================")
	fmt.Println("                        BADANIA ZAKOŃCZONE")
	fmt.Println("========================================================================")
}
