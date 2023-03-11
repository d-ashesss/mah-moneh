package api

import (
	"github.com/d-ashesss/mah-moneh/internal/accounts"
	"github.com/d-ashesss/mah-moneh/internal/categories"
	"github.com/d-ashesss/mah-moneh/internal/datastore"
	"github.com/d-ashesss/mah-moneh/internal/spendings"
	"github.com/d-ashesss/mah-moneh/internal/transactions"
)

var (
	ErrResourceNotFound = datastore.ErrRecordNotFound
)

type Service struct {
	accounts     *accounts.Service
	categories   *categories.Service
	transactions *transactions.Service
	spendings    *spendings.Service
}

func NewService(
	accounts *accounts.Service,
	categories *categories.Service,
	transactions *transactions.Service,
	spendings *spendings.Service,
) *Service {
	return &Service{
		accounts:     accounts,
		categories:   categories,
		transactions: transactions,
		spendings:    spendings,
	}
}
