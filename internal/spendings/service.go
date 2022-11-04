package spendings

import (
	"context"
	"github.com/d-ashesss/mah-moneh/internal/accounts"
	"github.com/d-ashesss/mah-moneh/internal/capital"
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

// Service is a service responsible for calculating spendings.
type Service struct {
	capital      CapitalService
	transactions TransactionsService
}

// NewService initializes the spendings service.
func NewService(cs CapitalService, ts TransactionsService) *Service {
	return &Service{capital: cs, transactions: ts}
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
	accounted, err := s.getTransactionSummary(ctx, u, month)
	if err != nil {
		return nil, err
	}
	totalSpent := s.subtractAmounts(currentCapital.Amounts, prevCapital.Amounts)
	sp := &Spending{
		Uncategorized: accounted,
		Unaccounted:   s.subtractAmounts(totalSpent, accounted),
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
func (s *Service) getTransactionSummary(ctx context.Context, u *users.User, month string) (map[string]float64, error) {
	sum := make(map[string]float64)
	txs, err := s.transactions.GetUserTransactions(ctx, u, month)
	if err != nil {
		return nil, err
	}
	for _, tx := range txs {
		sum[tx.Currency] += tx.Amount
	}
	return sum, nil
}

// subtractAmounts subtracts multi-currency amounts
func (s *Service) subtractAmounts(minuend, subtrahend map[string]float64) map[string]float64 {
	rest := make(map[string]float64)
	for currency, amount := range minuend {
		rest[currency] = amount
	}
	for currency, amount := range subtrahend {
		rest[currency] -= amount
	}
	return rest
}
