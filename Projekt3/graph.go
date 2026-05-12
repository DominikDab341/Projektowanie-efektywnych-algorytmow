package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type TSPInstance struct {
	Size   int
	Matrix [][]int
}

// ReadFromFile obsługuje zarówno pliki TSPLIB jak i proste z poprzedniego projektu
func ReadFromFile(filename string) (TSPInstance, error) {
	file, err := os.Open(filename)
	if err != nil {
		return TSPInstance{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	// Domyślny bufor (64 KB) jest za mały dla dużych plików TSPLIB (np. rbg443.atsp ~797 KB).
	// Zwiększamy do 1 MB, aby uniknąć błędu bufio.ErrTooLong.
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)
	scanner.Split(bufio.ScanLines)

	size := 0
	inMatrix := false
	var matrix [][]int
	currentRow := 0
	currentCol := 0

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if line == "EOF" {
			break
		}

		if inMatrix {
			words := strings.Fields(line)
			for _, word := range words {
				val, err := strconv.Atoi(word)
				if err != nil {
					return TSPInstance{}, fmt.Errorf("invalid matrix value: %v", word)
				}
				matrix[currentRow][currentCol] = val
				currentCol++
				if currentCol >= size {
					currentCol = 0
					currentRow++
				}
				if currentRow >= size {
					break
				}
			}
			if currentRow >= size {
				break
			}
			continue
		}

		// parsing headers
		parts := strings.SplitN(line, ":", 2)
		key := strings.TrimSpace(parts[0])

		if key == "DIMENSION" {
			if len(parts) > 1 {
				size, _ = strconv.Atoi(strings.TrimSpace(parts[1]))
			}
		} else if strings.HasPrefix(line, "DIMENSION") && size == 0 {
			// czasami nie ma dwukropka np. DIMENSION 171
			fields := strings.Fields(line)
			if len(fields) > 1 {
				size, _ = strconv.Atoi(fields[len(fields)-1])
			}
		}

		if line == "EDGE_WEIGHT_SECTION" {
			if size == 0 {
				return TSPInstance{}, fmt.Errorf("EDGE_WEIGHT_SECTION found before DIMENSION")
			}
			inMatrix = true
			matrix = make([][]int, size)
			for i := range matrix {
				matrix[i] = make([]int, size)
			}
		}
	}

	if !inMatrix {
		// Jeżeli to stary format bez TSPLIB nagłówków
		return ReadSimpleFile(filename)
	}

	return TSPInstance{Size: size, Matrix: matrix}, nil
}

// Fallback to Projekt 2 format if not TSPLIB
func ReadSimpleFile(filename string) (TSPInstance, error) {
	file, err := os.Open(filename)
	if err != nil {
		return TSPInstance{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	if !scanner.Scan() {
		return TSPInstance{}, fmt.Errorf("plik jest pusty lub uszkodzony")
	}
	size, _ := strconv.Atoi(scanner.Text())

	matrix := make([][]int, size)
	for i := range matrix {
		matrix[i] = make([]int, size)
	}

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

func (t TSPInstance) CalculatePathCost(path []int) int {
	if len(path) == 0 {
		return 0
	}
	cost := 0
	for i := 0; i < len(path)-1; i++ {
		cost += t.Matrix[path[i]][path[i+1]]
	}
	// wracamy do początku
	cost += t.Matrix[path[len(path)-1]][path[0]]
	return cost
}
