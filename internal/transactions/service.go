package transactions

import (
	"context"
	"github.com/d-ashesss/mah-moneh/internal/categories"
	"github.com/d-ashesss/mah-moneh/internal/users"
	"github.com/gofrs/uuid"
)

type Service struct {
	db Store
}

func NewService(db Store) *Service {
	return &Service{db: db}
}

func (s *Service) CreateTransaction(ctx context.Context, u *users.User, month string, currency string, amt float64, desc string, cat *categories.Category) (*Transaction, error) {
	tx := NewTransaction(u, month, currency, amt, desc, cat)
	if err := s.db.SaveTransaction(ctx, tx); err != nil {
		return nil, err
	}
	return tx, nil
}

func (s *Service) DeleteTransaction(ctx context.Context, tx *Transaction) error {
	return s.db.DeleteTransaction(ctx, tx)
}

func (s *Service) GetTransaction(ctx context.Context, uuid uuid.UUID) (*Transaction, error) {
	return s.db.GetTransaction(ctx, uuid)
}

func (s *Service) GetUserTransactions(ctx context.Context, u *users.User, month string) (TransactionCollection, error) {
	return s.db.GetUserTransactions(ctx, u, month)
}
