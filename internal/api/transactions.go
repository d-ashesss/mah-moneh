package api

import (
	"context"
	"github.com/d-ashesss/mah-moneh/internal/accounts"
	"github.com/d-ashesss/mah-moneh/internal/categories"
	"github.com/d-ashesss/mah-moneh/internal/transactions"
	"github.com/d-ashesss/mah-moneh/internal/users"
	"github.com/gofrs/uuid"
)

func (s *Service) CreateTransaction(
	ctx context.Context,
	u *users.User,
	month string,
	currency accounts.Currency,
	amt float64,
	desc string,
	cat *categories.Category,
) (*transactions.Transaction, error) {
	return s.transactions.CreateTransaction(ctx, u, month, currency, amt, desc, cat)
}

func (s *Service) DeleteTransaction(ctx context.Context, tx *transactions.Transaction) error {
	return s.transactions.DeleteTransaction(ctx, tx)
}

func (s *Service) GetTransaction(ctx context.Context, uuid uuid.UUID) (*transactions.Transaction, error) {
	return s.transactions.GetTransaction(ctx, uuid)
}

func (s *Service) GetUserTransactions(ctx context.Context, u *users.User, month string) (transactions.TransactionCollection, error) {
	return s.transactions.GetUserTransactions(ctx, u, month)
}
