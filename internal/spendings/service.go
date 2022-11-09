package spendings

import (
	"context"
	"github.com/d-ashesss/mah-moneh/internal/accounts"
	"github.com/d-ashesss/mah-moneh/internal/capital"
	"github.com/d-ashesss/mah-moneh/internal/categories"
	"github.com/d-ashesss/mah-moneh/internal/transactions"
	"github.com/d-ashesss/mah-moneh/internal/users"
	"time"
)

type CapitalService interface {
	GetCapital(ctx context.Context, u *users.User, month string) (*capital.Capital, error)
}

type TransactionsService interface {
	GetUserTransactions(ctx context.Context, u *users.User, month string) (transactions.TransactionCollection, error)
}

type CategoryService interface {
	GetUserCategories(ctx context.Context, u *users.User) ([]*categories.Category, error)
}

// Service is a service responsible for calculating spendings.
type Service struct {
	capital      CapitalService
	transactions TransactionsService
	categories   CategoryService
}

// NewService initializes the spendings service.
func NewService(capSrv CapitalService, transSrv TransactionsService, catSrv CategoryService) *Service {
	return &Service{capital: capSrv, transactions: transSrv, categories: catSrv}
}

// GetMonthSpendings calculates funds spent during specified month.
func (s *Service) GetMonthSpendings(ctx context.Context, u *users.User, month string) (*Spending, error) {
	currentCapital, err := s.capital.GetCapital(ctx, u, month)
	if err != nil {
		return nil, err
	}
	prevMonth, err := getPrevMonth(month)
	if err != nil {
		return nil, err
	}
	prevCapital, err := s.capital.GetCapital(ctx, u, prevMonth)
	if err != nil {
		return nil, err
	}
	categorized, uncategorized, err := s.getTransactionSummary(ctx, u, month)
	if err != nil {
		return nil, err
	}
	totalSpent := subtractAmounts(currentCapital.Amounts, prevCapital.Amounts)
	unaccounted := subtractAmounts(totalSpent, uncategorized)
	for _, spent := range categorized {
		unaccounted = subtractAmounts(unaccounted, spent)
	}
	sp := &Spending{
		ByCategory:    categorized,
		Uncategorized: uncategorized,
		Unaccounted:   unaccounted,
	}
	return sp, nil
}

// getPrevMonth calculates YYYY-MM representation of month previous to the provided.
func getPrevMonth(month string) (string, error) {
	d, err := time.Parse(accounts.FmtYearMonth, month)
	if err != nil {
		return "", err
	}
	return d.AddDate(0, -1, 0).Format(accounts.FmtYearMonth), nil
}

// getTransactionSummary calculates the sum of transactions recorded during given month.
func (s *Service) getTransactionSummary(ctx context.Context, u *users.User, month string) (map[*categories.Category]map[string]float64, map[string]float64, error) {
	sumByCat := make(map[*categories.Category]map[string]float64)
	sum := make(map[string]float64)
	txs, err := s.transactions.GetUserTransactions(ctx, u, month)
	if err != nil {
		return nil, nil, err
	}
	cats, err := s.categories.GetUserCategories(ctx, u)
	if err != nil {
		return nil, nil, err
	}
	for _, cat := range cats {
		if _, ok := sumByCat[cat]; !ok {
			sumByCat[cat] = make(map[string]float64)
		}
	}
txloop:
	for _, tx := range txs {
		for _, cat := range cats {
			if isSubset(cat.Tags, tx.Tags) {
				sumByCat[cat][tx.Currency] += tx.Amount
				continue txloop
			}
		}
		sum[tx.Currency] += tx.Amount
	}
	return sumByCat, sum, nil
}

// subtractAmounts subtracts multi-currency amounts
func subtractAmounts(minuend, subtrahend map[string]float64) map[string]float64 {
	rest := make(map[string]float64)
	for currency, amount := range minuend {
		rest[currency] = amount
	}
	for currency, amount := range subtrahend {
		rest[currency] -= amount
	}
	return rest
}

// isSubset determines if a slice is a subset of another slice
func isSubset(subset, in []string) bool {
	set := make(map[string]bool)
	for _, str := range in {
		set[str] = true
	}
	for _, str := range subset {
		if _, found := set[str]; !found {
			return false
		}
	}
	return true
}
