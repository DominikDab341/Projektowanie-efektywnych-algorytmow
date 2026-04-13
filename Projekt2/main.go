package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	var instance TSPInstance
	reader := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("\n--- MENU ---")
		fmt.Println("1. Wczytaj dane z pliku")
		fmt.Println("2. Wygeneruj dane losowe")
		fmt.Println("3. Wyświetl macierz")
		fmt.Println("4. Uruchom algorytm Branch & Bound dla obecnych danych")
		fmt.Println("5. Uruchom automatyczne badania i zbiór statystyk")
		fmt.Println("6. Weryfikacja poprawności B&B (porównanie z Brute-Force)")
		fmt.Println("0. Wyjście")
		fmt.Print("Wybierz opcję: ")

		reader.Scan()
		opcja := reader.Text()

		switch opcja {
		case "1":
			fmt.Print("Podaj nazwę pliku (np. tsp_10.txt): ")
			reader.Scan()
			path := reader.Text()

			var err error
			instance, err = ReadFromFile(path)
			if err != nil {
				fmt.Printf("Błąd podczas wczytywania: %v\n", err)
			} else {
				fmt.Println("Dane wczytane pomyślnie.")
			}

		case "2":
			fmt.Print("Podaj rozmiar N (liczbę miast): ")
			reader.Scan()
			nStr := reader.Text()
			n, err := strconv.Atoi(nStr)
			if err != nil || n <= 0 {
				fmt.Println("Błąd: podano nieprawidłowy rozmiar!")
			} else {
				instance = GenerateRandom(n)
				fmt.Printf("Wygenerowano losową macierz o rozmiarze %d.\n", n)
			}

		case "3":
			instance.Display()

		case "4":
			if instance.Size == 0 {
				fmt.Println("Błąd: Najpierw wczytaj lub wygeneruj dane wejściowe!")
				continue
			}

			fmt.Println("Jakim wariantem rozwiązania chcesz obrobić dane B&B?")
			fmt.Println("1) BEST (Best-First-Search ze stertą priorytetową)")
			fmt.Println("2) BREADTH (Breadth-First-Search po szerokości drzewa)")
			fmt.Print("Twój wybór: ")
			reader.Scan()
			
			wariant := reader.Text()
			metoda := "BEST"
			if wariant == "2" {
				metoda = "BREADTH"
			}

			fmt.Println("Jakie początkowe Ograniczenie wyznaczyć dla Pruningu gałęzi?")
			fmt.Println("1) Oparte o Nieskończoność (Startowe badanie w pełni ślepe)")
			fmt.Println("2) Oparte o wyznaczoną szybko drogę z Heurystyki (Najbliższy Sąsiad)")
			fmt.Print("Twój wybór: ")
			reader.Scan()

			modeOption := reader.Text()
			mode := "INF"
			if modeOption == "2" {
				mode = "NN"
			}

			fmt.Printf("Trwają obliczenia algorytmem Branch & Bound (%s / Ograniczenie początkowe: %s)...\n", metoda, mode)
			
			// W ręcznym sprawdzaniu zakładamy, że chcemy poczekać na wynik. Limit 0
			res := instance.SolveBranchAndBound(metoda, mode, 0)

			fmt.Printf("\n--- WYNIKI BRANCH & BOUND (%s) ---\n", metoda)
			fmt.Printf("Koszt najkrótszej trasy: %d\n", res.MinCost)
			
			fmt.Print("Ścieżka: ")
			for i, city := range res.Path {
				fmt.Print(city)
				if i < len(res.Path)-1 {
					fmt.Print(" -> ")
				}
			}
			fmt.Printf("\nCzas wykonania: %v\n", res.Duration)


		case "5":
			RunAutomatedTests()

		case "6":
			if instance.Size == 0 {
				fmt.Println("Błąd: Najpierw wczytaj lub wygeneruj dane wejściowe!")
				continue
			}
			if instance.Size > 12 {
				fmt.Println("Uwaga: Brute-force dla N>12 może trwać bardzo długo!")
			}

			fmt.Println("\n--- WERYFIKACJA POPRAWNOŚCI B&B vs BRUTE-FORCE ---")

			resBF := instance.SolveBruteForce()
			fmt.Printf("Brute-Force:    koszt = %d\n", resBF.MinCost)

			resBestINF := instance.SolveBranchAndBound("BEST", "INF", 0)
			fmt.Printf("Best(INF):      koszt = %d", resBestINF.MinCost)
			if resBestINF.MinCost == resBF.MinCost {
				fmt.Println("  ✓ OK")
			} else {
				fmt.Println("  ✗ BŁĄD!")
			}

			resBestNN := instance.SolveBranchAndBound("BEST", "NN", 0)
			fmt.Printf("Best(NN):       koszt = %d", resBestNN.MinCost)
			if resBestNN.MinCost == resBF.MinCost {
				fmt.Println("  ✓ OK")
			} else {
				fmt.Println("  ✗ BŁĄD!")
			}

			resBFS := instance.SolveBranchAndBound("BREADTH", "INF", 0)
			fmt.Printf("Breadth(INF):   koszt = %d", resBFS.MinCost)
			if resBFS.MinCost == resBF.MinCost {
				fmt.Println("  ✓ OK")
			} else {
				fmt.Println("  ✗ BŁĄD!")
			}

		case "0":
			fmt.Println("Zamykanie programu.")
			return

		default:
			fmt.Println("Nieznana opcja. Spróbuj ponownie.")
		}
	}
}
