package spendings

import (
	"github.com/d-ashesss/mah-moneh/internal/accounts"
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
	GetAmounts(cat *categories.Category) accounts.CurrencyAmounts
	AddTransaction(tx *transactions.Transaction)
}

// spendings contains calculated funds changes for a specific period of time.
type spendings map[*categories.Category]accounts.CurrencyAmounts

// newSpendings initializes new spendings structure.
func newSpendings(cats []*categories.Category) spendings {
	spent := make(spendings)
	for _, cat := range cats {
		if _, ok := spent[cat]; !ok {
			spent[cat] = accounts.NewCurrencyAmounts()
		}
	}
	spent[Uncategorized] = accounts.NewCurrencyAmounts()
	spent[Unaccounted] = accounts.NewCurrencyAmounts()
	spent[Total] = accounts.NewCurrencyAmounts()
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

func (s spendings) GetAmounts(cat *categories.Category) accounts.CurrencyAmounts {
	if amounts, found := s[cat]; found && amounts != nil {
		return amounts
	}
	return accounts.NewCurrencyAmounts()
}

func (s spendings) AddTransaction(tx *transactions.Transaction) {
	s.AddAmount(tx.Category, tx.Currency, tx.Amount)
}
