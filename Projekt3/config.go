package main

type CoolingScheme int

const (
	Geometric CoolingScheme = iota
	Linear
	LundyMees
)

type InitSolutionType int

const (
	RandomInit InitSolutionType = iota
	GreedyInit
)

type NeighborType int

const (
	Swap NeighborType = iota
	Invert
)

type SAConfig struct {
	InitialTemp  float64
	CoolingRate  float64 // alpha dla goemetric, beta dla LundyMees, zmniejszenie o stałą dla linear
	EpochLength  int
	MaxTimeMs    int // kryterium stopu w milisekundach
	Cooling      CoolingScheme
	InitSol      InitSolutionType
	NeighborGen  NeighborType
}

func DefaultConfig(size int) SAConfig {
	return SAConfig{
		InitialTemp: 10000.0,
		CoolingRate: 0.99,
		EpochLength: size,
		MaxTimeMs:   120000, // 2 minuty domyślnie
		Cooling:     Geometric,
		InitSol:     RandomInit,
		NeighborGen: Invert,
	}
}
