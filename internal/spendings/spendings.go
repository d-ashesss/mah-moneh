package spendings

import (
	"github.com/d-ashesss/mah-moneh/internal/accounts"
	"github.com/d-ashesss/mah-moneh/internal/categories"
	"github.com/d-ashesss/mah-moneh/internal/datastore"
	"github.com/d-ashesss/mah-moneh/internal/transactions"
	"github.com/gofrs/uuid"
)

// Service categories.
var (
	Uncategorized = &categories.Category{
		Model: datastore.Model{UUID: uuid.NewV5(uuid.Nil, "Uncategorized")},
	}
	Unaccounted = &categories.Category{
		Model: datastore.Model{UUID: uuid.NewV5(uuid.Nil, "Unaccounted")},
	}
	Total = &categories.Category{
		Model: datastore.Model{UUID: uuid.NewV5(uuid.Nil, "Total")},
	}
)

type Spendings interface {
	AddAmount(cat *categories.Category, currency string, amount float64)
	AddAmounts(cat *categories.Category, amounts accounts.CurrencyAmounts)
	GetAmount(cat *categories.Category, currency string) float64
	GetAmounts(cat *categories.Category) accounts.CurrencyAmounts
	AddTransaction(tx *transactions.Transaction)
}

// spendings contains calculated funds changes for a specific period of time.
type spendings map[uuid.UUID]accounts.CurrencyAmounts

// NewSpendings initializes new spendings structure.
func NewSpendings(cats []*categories.Category) Spendings {
	spent := make(spendings)
	for _, cat := range cats {
		if _, ok := spent[cat.UUID]; !ok {
			spent[cat.UUID] = accounts.NewCurrencyAmounts()
		}
	}
	spent[Uncategorized.UUID] = accounts.NewCurrencyAmounts()
	spent[Unaccounted.UUID] = accounts.NewCurrencyAmounts()
	spent[Total.UUID] = accounts.NewCurrencyAmounts()
	return spent
}

func (s spendings) AddAmount(cat *categories.Category, currency string, amount float64) {
	UUID := uuid.Nil
	if cat != nil {
		UUID = cat.UUID
	}
	if _, found := s[UUID]; found {
		s[UUID][currency] += amount
	} else {
		s[Uncategorized.UUID][currency] += amount
	}
	s[Total.UUID][currency] += amount
}

func (s spendings) AddAmounts(cat *categories.Category, amounts accounts.CurrencyAmounts) {
	for currency, amount := range amounts {
		s.AddAmount(cat, currency, amount)
	}
}

func (s spendings) GetAmount(cat *categories.Category, currency string) float64 {
	if amounts, found := s[cat.UUID]; found && amounts != nil {
		return amounts[currency]
	}
	return 0
}

func (s spendings) GetAmounts(cat *categories.Category) accounts.CurrencyAmounts {
	if amounts, found := s[cat.UUID]; found && amounts != nil {
		return amounts
	}
	return accounts.NewCurrencyAmounts()
}

func (s spendings) AddTransaction(tx *transactions.Transaction) {
	s.AddAmount(tx.Category, tx.Currency, tx.Amount)
}
