package spendings

import "github.com/d-ashesss/mah-moneh/internal/categories"

// Spending contains calculated funds changes for a specific period of time.
type Spending struct {
	ByCategory    map[*categories.Category]map[string]float64
	Uncategorized map[string]float64
	Unaccounted   map[string]float64
}

// New initializes new spending structure.
func New() *Spending {
	return &Spending{
		ByCategory:    make(map[*categories.Category]map[string]float64),
		Uncategorized: make(map[string]float64),
		Unaccounted:   make(map[string]float64),
	}
}
