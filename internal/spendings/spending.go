package spendings

// Spending contains calculated funds changes for a specific period of time.
type Spending struct {
	Uncategorized map[string]float64
	Unaccounted   map[string]float64
}

// New initializes new spending structure.
func New() *Spending {
	return &Spending{
		Uncategorized: make(map[string]float64),
		Unaccounted:   make(map[string]float64),
	}
}
