package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type TSPInstance struct {
	Size   int
	Matrix [][]int
}

func ReadFromFile(filename string) (TSPInstance, error) {
	file, err := os.Open(filename)
	if err != nil {
		return TSPInstance{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords) // Ignoruje spacje i entery, czyta tylko liczby

	if !scanner.Scan() {
		return TSPInstance{}, fmt.Errorf("plik jest pusty lub uszkodzony")
	}
	size, _ := strconv.Atoi(scanner.Text())

	// 2. Alokacja dynamiczna macierzy o rozmiarze size x size
	matrix := make([][]int, size)
	for i := range matrix {
		matrix[i] = make([]int, size)
	}

	// 3. Wczytywanie wartości macierzy
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			if scanner.Scan() {
				val, _ := strconv.Atoi(scanner.Text())
				matrix[i][j] = val
			}
		}
	}

	return TSPInstance{Size: size, Matrix: matrix}, nil
}


func (t TSPInstance) Display() {
	if t.Size == 0 {
		fmt.Println("Błąd: Macierz jest pusta. Najpierw wczytaj dane.")
		return
	}

	fmt.Printf("\n--- Macierz Kosztów (Rozmiar: %d) ---\n", t.Size)
	for i := 0; i < t.Size; i++ {
		for j := 0; j < t.Size; j++ {
			fmt.Printf("%4d ", t.Matrix[i][j])
		}
		fmt.Println()
	}
}

func main() {
	var instance TSPInstance
	reader := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("\n1. Wczytaj z pliku")
		fmt.Println("2. Wyświetl macierz")
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
			instance.Display()

		case "0":
			fmt.Println("Zamykanie programu.")
			return

		default:
			fmt.Println("Nieznana opcja.")
		}
	}
}