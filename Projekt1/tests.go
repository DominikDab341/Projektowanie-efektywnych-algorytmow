package main

import (
	"fmt"
)

// CalculateRelativeError zwraca błąd względny w procentach
func CalculateRelativeError(costAlg int, costBF int) float64 {
	if costBF == 0 {
		return 0.0 // Zabezpieczenie na wypadek dzielenia przez zero
	}
	
	roznica := float64(costAlg - costBF)
	return (roznica / float64(costBF)) * 100.0
}

// RunAutomatedTests wykonuje pełny zestaw testów wymaganych do sprawozdania
func RunAutomatedTests() {
	fmt.Println("\n=== ROZPOCZĘCIE TESTÓW AUTOMATYCZNYCH ===")
	
	// --- CZĘŚĆ 1: Pomiary czasu dla Brute-Force ---
	fmt.Println("\n1. Testy czasu wykonania Brute-Force [szukanie N ~ 2 min]")
	// Zaczynamy od mniejszych wartości, abyś mógł zaobserwować wzrost czasu
	testSizesBF := []int{8, 9, 10, 11, 12, 13} 
	for _, n := range testSizesBF {
		inst := GenerateRandom(n)
		res := inst.SolveBruteForce()
		fmt.Printf("N = %2d | Czas BF: %v\n", n, res.Duration)
	}

	// --- CZĘŚĆ 2: Pomiary błędu względnego dla NN, RNN, Random ---
	fmt.Println("\n2. Testy jakości algorytmów (błąd względny średni)")
	fmt.Println("N\tNN (%)\tRNN (%)\tRandom (%)")
	
	// Rozmiary wymagane w projekcie
	testSizesQuality := []int{10, 11, 12, 13, 14}
	
	for _, n := range testSizesQuality {
		var sumErrNN, sumErrRNN, sumErrRand float64
		
		fmt.Printf("Trwają obliczenia dla N = %d (100 iteracji). Proszę czekać...\n", n)
		
		for i := 0; i < 100; i++ {
			// Generujemy 100 macierzy odległości dla danego N
			inst := GenerateRandom(n)
			
			// 1. Uruchamiamy BF jako punkt odniesienia (rozwiązanie optymalne)
			resBF := inst.SolveBruteForce()
			optCost := resBF.MinCost
			
			// 2. Uruchamiamy NN
			resNN := inst.SolveNN(0) // Startujemy domyślnie od 0
			sumErrNN += CalculateRelativeError(resNN.MinCost, optCost)
			
			// 3. Uruchamiamy RNN
			resRNN := inst.SolveRNN()
			sumErrRNN += CalculateRelativeError(resRNN.MinCost, optCost)
			
			// 4. Uruchamiamy Random z ilością permutacji = 10 * N
			resRand := inst.SolveRandom(10 * n)
			sumErrRand += CalculateRelativeError(resRand.MinCost, optCost)
		}
		
		// Wyliczamy błąd średni
		avgErrNN := sumErrNN / 100.0
		avgErrRNN := sumErrRNN / 100.0
		avgErrRand := sumErrRand / 100.0
		
		// Wypisujemy wyniki w formie gotowej do wklejenia w Excela/tabelę
		fmt.Printf("%d\t%.2f\t%.2f\t%.2f\n", n, avgErrNN, avgErrRNN, avgErrRand)
	}
	fmt.Println("\n=== KONIEC TESTÓW ===")
}