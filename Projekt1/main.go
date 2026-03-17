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
		fmt.Println("\n--- MENU PROJEKTU PEA ---")
		fmt.Println("1. Wczytaj dane z pliku")
		fmt.Println("2. Wygeneruj dane losowe")
		fmt.Println("3. Wyświetl macierz")
		fmt.Println("4. Uruchom algorytm Brute-Force")
		fmt.Println("5. Uruchom algorytm Nearest Neighbor")
		fmt.Println("6. Uruchom algorytm Repetitive Nearest Neighbor")
		fmt.Println("7. Uruchom algorytm losowy")
		fmt.Println("8. Uruchom testy automatyczne")
		fmt.Println("0. Wyjście")
		fmt.Print("Wybierz opcję: ")

		reader.Scan()
		opcja := reader.Text()

		switch opcja {
		case "1":
			fmt.Print("Podaj nazwę pliku (np. dane.txt): ")
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
			fmt.Println("Trwają obliczenia Brute-Force...")
			res := instance.SolveBruteForce()

			fmt.Printf("\n--- WYNIKI BRUTE-FORCE ---\n")
			fmt.Printf("Koszt najkrótszej trasy: %d\n", res.MinCost)
			
			fmt.Print("Ścieżka: ")
			for i, city := range res.Path {
				fmt.Print(city)
				if i < len(res.Path)-1 {
					fmt.Print(" -> ")
				}
			}
			fmt.Printf("\nCzas wykonania: %v\n", res.Duration)

		case "0":
			fmt.Println("Zamykanie programu.")
			return

		case "5":
			if instance.Size == 0 {
				fmt.Println("Błąd: Najpierw wczytaj lub wygeneruj dane wejściowe!")
				continue
			}
			fmt.Println("Trwają obliczenia Nearest Neighbor...")
			
			res := instance.SolveNN(0)

			fmt.Printf("\n--- WYNIKI NEAREST NEIGHBOR ---\n")
			fmt.Printf("Koszt trasy: %d\n", res.MinCost)
			
			fmt.Print("Ścieżka: ")
			for i, city := range res.Path {
				fmt.Print(city)
				if i < len(res.Path)-1 {
					fmt.Print(" -> ")
				}
			}
			fmt.Printf("\nCzas wykonania: %v\n", res.Duration)

		case "6":
			if instance.Size == 0 {
				fmt.Println("Błąd: Najpierw wczytaj lub wygeneruj dane wejściowe!")
				continue
			}
			fmt.Println("Trwają obliczenia Repetitive Nearest Neighbor...")
			
			res := instance.SolveRNN()

			fmt.Printf("\n--- WYNIKI REPETITIVE NEAREST NEIGHBOR ---\n")
			fmt.Printf("Najmniejszy koszt trasy: %d\n", res.MinCost)
			
			fmt.Print("Ścieżka: ")
			for i, city := range res.Path {
				fmt.Print(city)
				if i < len(res.Path)-1 {
					fmt.Print(" -> ")
				}
			}
			fmt.Printf("\nCzas wykonania: %v\n", res.Duration)


		case "7":
			if instance.Size == 0 {
				fmt.Println("Błąd: Najpierw wczytaj lub wygeneruj dane wejściowe!")
				continue
			}
			fmt.Print("Podaj liczbę permutacji do wylosowania: ")
			reader.Scan()
			perms, err := strconv.Atoi(reader.Text())
			if err != nil || perms <= 0 {
				fmt.Println("Błąd: podano nieprawidłową liczbę permutacji!")
				continue
			}

			fmt.Printf("Trwają obliczenia dla %d losowych permutacji...\n", perms)
			res := instance.SolveRandom(perms)

			fmt.Printf("\n--- WYNIKI ALGORYTMU LOSOWEGO ---\n")
			fmt.Printf("Najmniejszy znaleziony koszt: %d\n", res.MinCost)
			
			fmt.Print("Ścieżka: ")
			for i, city := range res.Path {
				fmt.Print(city)
				if i < len(res.Path)-1 {
					fmt.Print(" -> ")
				}
			}
			fmt.Printf("\nCzas wykonania: %v\n", res.Duration)

		case "8":
			RunAutomatedTests()

		default:
			fmt.Println("Nieznana opcja. Spróbuj ponownie.")
		}
	}
}