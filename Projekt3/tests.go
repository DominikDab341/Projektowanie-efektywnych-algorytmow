package main

import (
	"fmt"
	"strings"
)

var bestKnownSolutions = map[string]int{
	"br17.atsp":    39,
	"ftv33.atsp":   1286,
	"p43.atsp":     5620,
	"ft53.atsp":    6905,
	"ft70.atsp":    38673,
	"kro124p.atsp": 36230,
	"ftv170.atsp":  2755,
	"rbg323.atsp":  1326,
	"rbg403.atsp":  2465,
	"rbg443.atsp":  2720,
}

type TestResult struct {
	InstanceName string
	Parameter    string
	AvgCost      float64
	AvgTimeMs    float64
	PRD          float64
}

func RunAutomaticTests() {
	fmt.Println("--- ROZPOCZYNANIE TESTÓW AUTOMATYCZNYCH ---")
	
	// Lista 10 plików o mocno różniących się rozmiarach
	instances := []string{
		"br17.atsp", "ftv33.atsp", "p43.atsp", "ft53.atsp", 
		"ft70.atsp", "kro124p.atsp", "ftv170.atsp", "rbg323.atsp", 
		"rbg403.atsp", "rbg443.atsp",
	}
    
    // 1. TEST 3.0: Skalowalność (różne instancje)
    runTest30(instances)

	// 2. TEST 3.5: Wpływ rozwiązania początkowego
	runTest35(instances)

	// 3. TEST 4.0: Wpływ schematów chłodzenia
	runTest40(instances)
}

func runTest35(files []string) {
	fmt.Println("\n--- TEST 3.5: Rozwiązanie początkowe (10 prób na wariant, różne instancje) ---")
	printTableHeader()

	methods := []InitSolutionType{RandomInit, GreedyInit}
	methodNames := []string{"Losowe", "Zachłanne"}

	for _, f := range files {
		instance, err := ReadFromFile(f)
		if err != nil { continue }

		for i, method := range methods {
			var totalCost int
			var totalTime int64
			runs := 10

			for r := 0; r < runs; r++ {
				config := SAConfig{
					MaxTimeMs:   10000, // 10s na próbę
					EpochLength: instance.Size * 50,
					Cooling:     Geometric,
					CoolingRate: 0.99,
					InitSol:     method,
					NeighborGen: Swap,
				}
				sa := SimulatedAnnealing{Instance: instance, Config: config}
				sa.Config.InitialTemp = sa.CalculateInitialTemp(0.99, 500)
				
				res := sa.Solve()
				totalCost += res.MinCost
				totalTime += res.Duration.Milliseconds()
			}
			label := fmt.Sprintf("%s (%s)", f, methodNames[i])
			printRow(label, f, totalCost, totalTime, runs)
		}
	}
}

func runTest40(files []string) {
	fmt.Println("\n--- TEST 4.0: Schematy chłodzenia (5 prób na wariant, różne instancje) ---")
	printTableHeader()

	schemes := []struct {
		id   CoolingScheme
		name string
		rate float64
	}{
		{Geometric, "Geometric", 0.99},
		{LundyMees, "Lundy-Mees", 0.001},
	}

	for _, f := range files {
		instance, err := ReadFromFile(f)
		if err != nil { continue }

		for _, s := range schemes {
			var totalCost int
			var totalTime int64
			runs := 10

			for r := 0; r < runs; r++ {
				config := SAConfig{
					MaxTimeMs:   10000,
					EpochLength: instance.Size * 50,
					Cooling:     s.id,
					CoolingRate: s.rate,
					InitSol:     GreedyInit,
					NeighborGen: Swap,
				}
				sa := SimulatedAnnealing{Instance: instance, Config: config}
				sa.Config.InitialTemp = sa.CalculateInitialTemp(0.99, 500)
				
				res := sa.Solve()
				totalCost += res.MinCost
				totalTime += res.Duration.Milliseconds()
			}
			label := fmt.Sprintf("%s (%s)", f, s.name)
			printRow(label, f, totalCost, totalTime, runs)
		}
	}
}

func runTest30(files []string) {
    fmt.Println("\n--- TEST 3.0: Skalowalność (Czas i Błąd vs Rozmiar) ---")
    printTableHeader()

    for _, f := range files {
        instance, err := ReadFromFile(f)
        if err != nil { continue }
        
        var totalCost int
        var totalTime int64
        runs := 5 // mniej prób dla dużych instancji, żeby oszczędzić czas

        for r := 0; r < runs; r++ {
            config := SAConfig{
                MaxTimeMs:   60000, // 60s - dla bardzo dużych jak rbg403 może być potrzebne więcej (do 15 minut)
                EpochLength: instance.Size * 100,
                Cooling:     Geometric,
                CoolingRate: 0.995,
                InitSol:     GreedyInit,
                NeighborGen: Swap,
            }
            sa := SimulatedAnnealing{Instance: instance, Config: config}
            sa.Config.InitialTemp = sa.CalculateInitialTemp(0.99, 500)
            res := sa.Solve()
            totalCost += res.MinCost
            totalTime += res.Duration.Milliseconds()
        }
        printRow(fmt.Sprintf("%s (%d)", f, instance.Size), f, totalCost, totalTime, runs)
    }
}

func printRow(label string, filename string, totalCost int, totalTime int64, runs int) {
	avgCost := float64(totalCost) / float64(runs)
	avgTime := float64(totalTime) / float64(runs)
	
	bks := bestKnownSolutions[filename]
	prd := 0.0
	if bks > 0 {
		prd = ((avgCost - float64(bks)) / float64(bks)) * 100
	}
	fmt.Printf("%-30s | %15.2f | %15.2f | %9.2f%%\n", label, avgCost, avgTime, prd)
}

func printTableHeader() {
	fmt.Printf("%-30s | %15s | %15s | %10s\n", "Parametr/Instancja", "Średni Koszt", "Śr. Czas [ms]", "Błąd PRD")
	fmt.Println(strings.Repeat("-", 79))
}