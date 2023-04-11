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
func (s *Service) GetMonthSpendings(ctx context.Context, u *users.User, month string) (Spendings, error) {
	spent, err := s.getTransactionSummary(ctx, u, month)
	if err != nil {
		return nil, err
	}
	capt, err := s.getCapitalDiff(ctx, u, month)
	if err != nil {
		return nil, err
	}
	unaccounted := capt.Diff(spent.GetAmounts(Total))
	spent.AddAmounts(Unaccounted, unaccounted)
	return spent, nil
}

// getPrevMonth calculates YYYY-MM representation of month previous to the provided.
func getPrevMonth(month string) (string, error) {
	d, err := time.Parse(accounts.FmtYearMonth, month)
	if err != nil {
		return "", err
	}
	return d.AddDate(0, -1, 0).Format(accounts.FmtYearMonth), nil
}

// getCapitalDiff calculates the difference between specified month and previous month capitals.
func (s *Service) getCapitalDiff(ctx context.Context, u *users.User, month string) (accounts.CurrencyAmounts, error) {
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
	return currentCapital.Diff(prevCapital), nil
}

// getTransactionSummary calculates the sum of transactions recorded during given month.
func (s *Service) getTransactionSummary(ctx context.Context, u *users.User, month string) (Spendings, error) {
	cats, err := s.categories.GetUserCategories(ctx, u)
	if err != nil {
		return nil, err
	}
	spent := NewSpendings(cats)

	txs, err := s.transactions.GetUserTransactions(ctx, u, month)
	if err != nil {
		return nil, err
	}
	for _, tx := range txs {
		spent.AddTransaction(tx)
	}
	return spent, nil
}
