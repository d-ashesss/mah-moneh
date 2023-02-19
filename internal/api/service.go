package api

import (
	"github.com/d-ashesss/mah-moneh/internal/accounts"
	"github.com/d-ashesss/mah-moneh/internal/categories"
	"github.com/d-ashesss/mah-moneh/internal/transactions"
)

type Service struct {
	accounts     *accounts.Service
	categories   *categories.Service
	transactions *transactions.Service
}

func NewService(
	accounts *accounts.Service,
	categories *categories.Service,
	transactions *transactions.Service,
) *Service {
	return &Service{
		accounts:     accounts,
		categories:   categories,
		transactions: transactions,
	}
}
