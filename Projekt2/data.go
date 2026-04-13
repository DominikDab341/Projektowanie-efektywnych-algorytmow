package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

// Przechowuje rozmiar i dynamicznie alokowaną macierz
type TSPInstance struct {
	Size   int
	Matrix [][]int
}

// Przechowuje wyniki działania algorytmów
type Result struct {
	Path          []int
	MinCost       int
	Duration      time.Duration
	NodesExplored int
}

// wczytuje dane z pliku tekstowego
func ReadFromFile(filename string) (TSPInstance, error) {
	file, err := os.Open(filename)
	if err != nil {
		return TSPInstance{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	// Wczytywanie rozmiaru macierzy
	if !scanner.Scan() {
		return TSPInstance{}, fmt.Errorf("plik jest pusty lub uszkodzony")
	}
	size, _ := strconv.Atoi(scanner.Text())

	// Dynamiczna alokacja macierzy
	matrix := make([][]int, size)
	for i := range matrix {
		matrix[i] = make([]int, size)
	}

	// Wypełnianie macierzy wartościami
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

// Generuje losową instancję ATSP o zadanym rozmiarze
func GenerateRandom(size int) TSPInstance {
	matrix := make([][]int, size)

	for i := 0; i < size; i++ {
		matrix[i] = make([]int, size)
		for j := 0; j < size; j++ {
			if i == j {
				matrix[i][j] = -1
			} else {
				matrix[i][j] = rand.Intn(90) + 10 // Koszty z zakresu 10-99
			}
		}
	}
	
	return TSPInstance{Size: size, Matrix: matrix}
}

// Wypisuje zawartość macierzy w czytelny sposób
func (t TSPInstance) Display() {
	if t.Size == 0 {
		fmt.Println("Błąd: Macierz jest pusta. Najpierw wczytaj lub wygeneruj dane.")
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
