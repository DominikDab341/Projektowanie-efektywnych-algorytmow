package main

import (
	"fmt"
)

func CalculateRelativeError(costAlg int, costBF int) float64 {
	if costBF == 0 {
		return 0.0
	}
	
	roznica := float64(costAlg - costBF)
	return (roznica / float64(costBF)) * 100.0
}

func RunAutomatedTests() {
	// Pomiary czasu dla Brute-Force
	fmt.Println("\n1. Testy czasu wykonania Brute-Force")
	testSizesBF := []int{8, 9, 10, 11, 12, 13, 14} 
	for _, n := range testSizesBF {
		inst := GenerateRandom(n)
		res := inst.SolveBruteForce()
		fmt.Printf("N = %2d | Czas BF: %v\n", n, res.Duration)
	}

	fmt.Println("\n2. Testy jakości algorytmów (błąd względny średni)")
	fmt.Println("N\tNN (%)\tRNN (%)\tRandom (%)")
	
	testSizesQuality := []int{10, 11, 12, 13, 14}
	
	for _, n := range testSizesQuality {
		var sumErrNN, sumErrRNN, sumErrRand float64
		
		fmt.Printf("Trwają obliczenia dla N = %d. Proszę czekać...\n", n)
		
		for i := 0; i < 100; i++ {
			inst := GenerateRandom(n)
			
			resBF := inst.SolveBruteForce()
			optCost := resBF.MinCost
			
			resNN := inst.SolveNN(0)
			sumErrNN += CalculateRelativeError(resNN.MinCost, optCost)
			
			resRNN := inst.SolveRNN()
			sumErrRNN += CalculateRelativeError(resRNN.MinCost, optCost)
			
			resRand := inst.SolveRandom(10 * n)
			sumErrRand += CalculateRelativeError(resRand.MinCost, optCost)
		}
		
		// Wyliczamy błąd średni
		avgErrNN := sumErrNN / 100.0
		avgErrRNN := sumErrRNN / 100.0
		avgErrRand := sumErrRand / 100.0
		
		fmt.Printf("%d\t%.2f\t%.2f\t%.2f\n", n, avgErrNN, avgErrRNN, avgErrRand)
	}
}