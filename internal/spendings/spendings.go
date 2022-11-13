package spendings

import (
	"github.com/d-ashesss/mah-moneh/internal/categories"
	"github.com/d-ashesss/mah-moneh/internal/transactions"
)

// Service categories.
var (
	Uncategorized = &categories.Category{}
	Unaccounted   = &categories.Category{}
	Total         = &categories.Category{}
)

type Spendings interface {
	AddAmount(cat *categories.Category, currency string, amount float64)
	GetAmount(cat *categories.Category, currency string) float64
	AddTransaction(tx *transactions.Transaction)
}

// spendings contains calculated funds changes for a specific period of time.
type spendings map[*categories.Category]map[string]float64

// newSpendings initializes new spending structure.
func newSpendings(cats []*categories.Category) spendings {
	spent := make(spendings)
	for _, cat := range cats {
		if _, ok := spent[cat]; !ok {
			spent[cat] = make(map[string]float64)
		}
	}
	spent[Uncategorized] = make(map[string]float64)
	spent[Unaccounted] = make(map[string]float64)
	spent[Total] = make(map[string]float64)
	return spent
}

func (s spendings) AddAmount(cat *categories.Category, currency string, amount float64) {
	if _, found := s[cat]; found {
		s[cat][currency] += amount
	} else {
		s[Uncategorized][currency] += amount
	}
	s[Total][currency] += amount
}

func (s spendings) GetAmount(cat *categories.Category, currency string) float64 {
	if amounts, found := s[cat]; found && amounts != nil {
		return amounts[currency]
	}
	return 0
}

func (s spendings) AddTransaction(tx *transactions.Transaction) {
	s.AddAmount(tx.Category, tx.Currency, tx.Amount)
}
