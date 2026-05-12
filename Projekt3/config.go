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
	Insert
)

type SAConfig struct {
	InitialTemp float64
	CoolingRate float64 // alpha dla geometric, beta dla LundyMees, zmniejszenie o stałą dla linear
	EpochLength int
	MaxTimeMs   int // kryterium stopu w milisekundach
	Cooling     CoolingScheme
	InitSol     InitSolutionType
	NeighborGen NeighborType
}

func DefaultConfig(size int) SAConfig {
	return SAConfig{
		MaxTimeMs:   10000, // 10 sekund
		InitialTemp: 1000.0,
		EpochLength: size,
		Cooling:     Geometric,
		CoolingRate: 0.99,
		InitSol:     RandomInit,
		NeighborGen: Insert,
	}
}
