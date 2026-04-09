package main

import (
	"fmt"
	"time"
)

// RunAutomatedTests przeprowadza statystyczną analizę działania Branch&Bound ATSP wg wymogów zadania 2.
// Funkcja zakłada testy dla zróżnicowanych 15 rozmiarów (N), ze 100 próbkami z zachowaniem odcięcia 5 min.
func RunAutomatedTests() {
	LimitCzasu := 5 * time.Minute
	LimitInstancji := 100 // w zadaniu: 100 losowych instancji problemu dla każdej wielkości

	fmt.Println("\n--- ROZPOCZYNANIE BADAŃ B&B (Limit: 5 minut na przerwę pojedynczej) ---")
	
	// Zalecane minimum 10 rozmiarów. B&B wzrośnie drastycznie pomiędzy N=10 a N=20.
	testSizes := []int{4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	
	fmt.Println("N\tBFS (INF) T.\tPrzerw BFS%\tBest (INF) T.\tPrzerw Best%\tBest (NN) T.\tPrzerw Best(NN)%")

	for _, n := range testSizes {
		var (
			sumTimeBFS, sumTimeBestINF, sumTimeBestNN time.Duration
			timeoutBFS, timeoutBestINF, timeoutBestNN int
		)
		
		for i := 0; i < LimitInstancji; i++ {
			// Czas wygenerowania instancji losowej N i alokacji nie wpływa na pomiar funkcji B&B
			inst := GenerateRandom(n)
			
			// === TEST 1: Metoda BREADTH FIRST SEARCH (Ślepy start) ===
			resBFS := inst.SolveBranchAndBound("BREADTH", "INF", LimitCzasu)
			if resBFS.MinCost == -1 {
				timeoutBFS++ 
			} else {
				sumTimeBFS += resBFS.Duration
			}
			
			// === TEST 2: Metoda BEST FIRST SEARCH (Ślepy start) ===
			resBestINF := inst.SolveBranchAndBound("BEST", "INF", LimitCzasu)
			if resBestINF.MinCost == -1 {
				timeoutBestINF++
			} else {
				sumTimeBestINF += resBestINF.Duration
			}

			// === TEST 3: Metoda BEST FIRST SEARCH (Heurystyka NN na start) ===
			resBestNN := inst.SolveBranchAndBound("BEST", "NN", LimitCzasu)
			if resBestNN.MinCost == -1 {
				timeoutBestNN++
			} else {
				sumTimeBestNN += resBestNN.Duration
			}
		}

		// Obliczenie uśrednionych statystyk prób
		sukcesyBFS := LimitInstancji - timeoutBFS
		sukcesyBestINF := LimitInstancji - timeoutBestINF
		sukcesyBestNN := LimitInstancji - timeoutBestNN
		
		var avgBFS, avgBestINF, avgBestNN time.Duration
		
		if sukcesyBFS > 0 { avgBFS = time.Duration(float64(sumTimeBFS.Nanoseconds()) / float64(sukcesyBFS)) } else { avgBFS = LimitCzasu }
		if sukcesyBestINF > 0 { avgBestINF = time.Duration(float64(sumTimeBestINF.Nanoseconds()) / float64(sukcesyBestINF)) } else { avgBestINF = LimitCzasu }
		if sukcesyBestNN > 0 { avgBestNN = time.Duration(float64(sumTimeBestNN.Nanoseconds()) / float64(sukcesyBestNN)) } else { avgBestNN = LimitCzasu }

		// Z uwagi, że testów w grupie jest równo 100, licznik `timeout` jest wprost równoznaczny procencie (%).
		fmt.Printf("%d\t%v\t%d%%\t\t%v\t%d%%\t\t%v\t%d%%\n", 
			n, avgBFS, timeoutBFS, avgBestINF, timeoutBestINF, avgBestNN, timeoutBestNN)
		
		// Oszacowanie maksymalnego osiągalnego średniego N - jeżeli dla WSZYSTKICH wariantów padnie, przerywamy.
		if timeoutBFS == LimitInstancji && timeoutBestINF == LimitInstancji && timeoutBestNN == LimitInstancji {
			fmt.Printf("Osiągnięto blokadę obliczeniową przy N=%d, limitu pamięci/czasu. Zakańczam.\n", n)
			break
		}
	}
}
