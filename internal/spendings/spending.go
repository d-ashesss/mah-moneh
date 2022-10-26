package spendings

// Spending contains calculated funds changes for a specific period of time.
type Spending struct {
	Amounts map[string]float64
}

// New initializes new spending structure.
func New() *Spending {
	a := make(map[string]float64)
	return &Spending{Amounts: a}
}
