package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	var instance TSPInstance
	var config SAConfig = DefaultConfig(0)

	reader := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("\n--- MENU SYMULOWANE WYŻARZANIE (SA) ---")
		fmt.Println("1. Wczytanie danych z pliku (TSPLIB)")
		fmt.Printf("2. Wprowadzenie kryterium stopu (Obecnie: %d ms)\n", config.MaxTimeMs)
		fmt.Println("3. Obliczanie i wyświetlenie rozwiązania początkowego")
		fmt.Println("4. Modyfikacja domyślnych ustawień algorytmu")
		fmt.Println("5. Uruchomienie algorytmu SA")
		fmt.Println("6. URUCHOM TESTY AUTOMATYCZNE")
		fmt.Println("0. Wyjście")
		fmt.Print("Wybierz opcję: ")

		if !reader.Scan() {
			break
		}
		opcja := reader.Text()

		switch opcja {
		case "1":
			fmt.Print("Podaj nazwę pliku (np. ft53.atsp): ")
			reader.Scan()
			path := reader.Text()

			var err error
			instance, err = ReadFromFile(path)
			if err != nil {
				fmt.Printf("Błąd podczas wczytywania: %v\n", err)
			} else {
				fmt.Printf("Dane wczytane pomyślnie. Rozmiar: %d\n", instance.Size)
				config.EpochLength = instance.Size * 100
			}

		case "2":
			fmt.Print("Podaj kryterium stopu w sekundach: ")
			reader.Scan()
			s, _ := strconv.Atoi(reader.Text())
			if s > 0 {
				config.MaxTimeMs = s * 1000
				fmt.Printf("Limit czasu: %d s\n", s)
			}

		case "3":
			if instance.Size == 0 {
				fmt.Println("Najpierw wczytaj dane!")
				continue
			}
			sa := SimulatedAnnealing{Instance: instance, Config: config}
			initPath := sa.generateInitialSolution()
			fmt.Printf("\nKoszt początkowy (%s): %d\n", getInitTypeString(config.InitSol), instance.CalculatePathCost(initPath))
			fmt.Printf("Droga (%d węzłów, indeksy od 1): %s\n", len(initPath), formatPath(initPath))

		case "4":
			modyfikujUstawienia(&config, instance, reader)

		case "5":
			if instance.Size == 0 {
				fmt.Println("Błąd: Brak wczytanej instancji!")
				continue
			}
			sa := SimulatedAnnealing{Instance: instance, Config: config}
			fmt.Println("Optymalizacja w toku...")
			res := sa.Solve()
			fmt.Printf("\nZakończono! Najlepszy koszt: %d (Czas: %v)\n", res.MinCost, res.Duration)
			fmt.Printf("Droga (%d węzłów, indeksy od 1): %s\n", len(res.Path), formatPath(res.Path))

		case "6":
			RunAutomaticTests()

		case "0":
			fmt.Println("Zamykanie.")
			return

		default:
			fmt.Println("Nieznana opcja.")
		}
	}
}

// modyfikujUstawienia pozwala użytkownikowi zmienić parametry algorytmu.
// Przyjmuje wskaźnik na istniejący scanner (zamiast tworzyć drugi na os.Stdin).
func modyfikujUstawienia(config *SAConfig, inst TSPInstance, reader *bufio.Scanner) {

	fmt.Println("\n--- MODYFIKACJA USTAWIEŃ ---")
	fmt.Println("1. Chłodzenie: 1-Geo, 2-Lin, 3-Lundy")
	fmt.Print("Wybór: ")
	reader.Scan()
	switch reader.Text() {
	case "1":
		config.Cooling = Geometric
		config.CoolingRate = 0.99
	case "2":
		config.Cooling = Linear
		config.CoolingRate = 1.0
	case "3":
		config.Cooling = LundyMees
		config.CoolingRate = 0.0001
	}

	fmt.Println("2. Start: 1-Losowy, 2-Zachłanny")
	fmt.Print("Wybór: ")
	reader.Scan()
	if reader.Text() == "2" {
		config.InitSol = GreedyInit
	} else {
		config.InitSol = RandomInit
	}

	fmt.Println("3. Sąsiedztwo: 1-Swap, 2-Invert")
	fmt.Print("Wybór: ")
	reader.Scan()
	if reader.Text() == "2" {
		config.NeighborGen = Invert
	} else {
		config.NeighborGen = Swap
	}

	fmt.Print("Czy przeliczyć T0 automatycznie? (T/N): ")
	reader.Scan()
	if strings.ToUpper(reader.Text()) == "T" && inst.Size > 0 {
		sa := SimulatedAnnealing{Instance: inst, Config: *config}
		config.InitialTemp = sa.CalculateInitialTemp(0.99, 1000)
		fmt.Printf("Nowe T0: %.2f\n", config.InitialTemp)
	}
}

func getInitTypeString(it InitSolutionType) string {
	if it == GreedyInit {
		return "Zachłanne"
	}
	return "Losowe"
}

// formatPath formatuje ścieżkę jako string z węzłami numerowanymi od 1
// zgodnie z konwencją TSPLIB (wewnętrznie przechowujemy od 0).
func formatPath(path []int) string {
	var sb strings.Builder
	sb.WriteByte('[')
	for i, node := range path {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(node + 1))
	}
	sb.WriteByte(']')
	return sb.String()
}
