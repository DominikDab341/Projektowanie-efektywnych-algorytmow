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
			
			// Wywołanie funkcji z pliku bruteforce.go
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

		default:
			fmt.Println("Nieznana opcja. Spróbuj ponownie.")
		}
	}
}