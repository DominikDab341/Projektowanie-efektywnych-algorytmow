package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	var instance TSPInstance
	var config SAConfig = DefaultConfig(0) // 0 rozmiar na razie

	reader := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("\n--- MENU SYMULOWANE WYŻARZANIE (SA) ---")
		fmt.Println("1. Wczytanie danych z pliku (TSPLIB)")
		fmt.Printf("2. Wprowadzenie kryterium stopu (Obecnie: %d ms)\n", config.MaxTimeMs)
		fmt.Println("3. Obliczanie i wyświetlenie rozwiązania początkowego")
		fmt.Println("4. Modyfikacja domyślnych ustawień algorytmu (Chłodzenie, Sąsiedztwo, Temperatura)")
		fmt.Println("5. Uruchomienie algorytmu SA")
		fmt.Println("0. Wyjście")
		fmt.Print("Wybierz opcję: ")

		if !reader.Scan() {
			fmt.Println("\nZakończono wejście (EOF). Zamykanie.")
			break
		}
		opcja := reader.Text()

		switch opcja {
		case "1":
			fmt.Print("Podaj nazwę pliku (np. ftv170.atsp): ")
			reader.Scan()
			path := reader.Text()

			var err error
			instance, err = ReadFromFile(path)
			if err != nil {
				fmt.Printf("Błąd podczas wczytywania: %v\n", err)
			} else {
				fmt.Printf("Dane wczytane pomyślnie. Rozmiar macierzy: %d\n", instance.Size)
				// Aktualizacja długości epoki po wczytaniu rozmiaru
				config.EpochLength = instance.Size
			}

		case "2":
			fmt.Print("Podaj kryterium stopu w sekundach: ")
			reader.Scan()
			sStr := reader.Text()
			s, err := strconv.Atoi(sStr)
			if err != nil || s <= 0 {
				fmt.Println("Błąd: podano nieprawidłowy czas!")
			} else {
				config.MaxTimeMs = s * 1000
				fmt.Printf("Ustawiono limit czasu na %d sekund.\n", s)
			}

		case "3":
			if instance.Size == 0 {
				fmt.Println("Błąd: Najpierw wczytaj dane wejściowe!")
				continue
			}
			
			sa := SimulatedAnnealing{Instance: instance, Config: config}
			initPath := sa.generateInitialSolution()
			initCost := instance.CalculatePathCost(initPath)
			
			fmt.Printf("\n--- Rozwiązanie początkowe (%s) ---\n", getInitTypeString(config.InitSol))
			fmt.Printf("Koszt trasy: %d\n", initCost)
			fmt.Print("Pierwsze 10 miast ścieżki: ")
			limit := 10
			if len(initPath) < limit {
				limit = len(initPath)
			}
			for i := 0; i < limit; i++ {
				fmt.Printf("%d ", initPath[i])
			}
			fmt.Println("...")

		case "4":
			fmt.Println("\n--- MODYFIKACJA USTAWIEŃ ---")
			fmt.Println("Wybierz schemat chłodzenia:")
			fmt.Println("1. Geometryczne (T = T * alpha)")
			fmt.Println("2. Liniowe (T = T - delta)")
			fmt.Println("3. Lundy-Mees (T = T / (1 + beta*T))")
			fmt.Print("Wybór: ")
			reader.Scan()
			switch reader.Text() {
			case "1": config.Cooling = Geometric; config.CoolingRate = 0.9999
			case "2": config.Cooling = Linear; config.CoolingRate = 0.1
			case "3": config.Cooling = LundyMees; config.CoolingRate = 0.001
			default: fmt.Println("Nieznana opcja, pozostawiono domyślne.")
			}

			fmt.Println("\nWybierz rozwiązanie początkowe:")
			fmt.Println("1. Losowe")
			fmt.Println("2. Zachłanne (Najbliższy Sąsiad)")
			fmt.Print("Wybór: ")
			reader.Scan()
			if reader.Text() == "2" {
				config.InitSol = GreedyInit
			} else {
				config.InitSol = RandomInit
			}

			fmt.Println("\nWybierz sąsiedztwo:")
			fmt.Println("1. Swap (zamiana dwóch miast)")
			fmt.Println("2. Invert (odwrócenie podciągu)")
			fmt.Print("Wybór: ")
			reader.Scan()
			if reader.Text() == "1" {
				config.NeighborGen = Swap
			} else {
				config.NeighborGen = Invert
			}

			fmt.Println("\nCzy chcesz przeliczyć temperaturę początkową automatycznie na podstawie instancji? (T/N)")
			fmt.Print("Wybór: ")
			reader.Scan()
			if reader.Text() == "T" || reader.Text() == "t" {
				if instance.Size > 0 {
					sa := SimulatedAnnealing{Instance: instance, Config: config}
					newTemp := sa.CalculateInitialTemp(0.99, 1000)
					config.InitialTemp = newTemp
					fmt.Printf("Nowa temperatura początkowa (P=0.99): %.2f\n", newTemp)
				} else {
					fmt.Println("Błąd: Wczytaj najpierw instancję!")
				}
			}

		case "5":
			if instance.Size == 0 {
				fmt.Println("Błąd: Najpierw wczytaj dane wejściowe!")
				continue
			}

			fmt.Println("\nUruchamiam algorytm Symulowanego Wyżarzania (SA)...")
			sa := SimulatedAnnealing{Instance: instance, Config: config}
			res := sa.Solve()

			fmt.Printf("\n--- WYNIKI SYMULOWANEGO WYŻARZANIA ---\n")
			fmt.Printf("Koszt optymalnej znalezionej trasy: %d\n", res.MinCost)
			fmt.Printf("Czas wykonania: %v\n", res.Duration)
			
			fmt.Println("\nŚcieżka (może być bardzo długa):")
			for i, city := range res.Path {
				fmt.Printf("%d ", city)
				if i < len(res.Path)-1 && (i+1)%20 == 0 {
					fmt.Println()
				}
			}
			fmt.Println()

		case "0":
			fmt.Println("Zamykanie programu.")
			return

		default:
			fmt.Println("Nieznana opcja. Spróbuj ponownie.")
		}
	}
}

func getInitTypeString(it InitSolutionType) string {
	if it == GreedyInit {
		return "Zachłanne / NN"
	}
	return "Losowe"
}

