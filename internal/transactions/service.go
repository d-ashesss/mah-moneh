package transactions

import (
	"context"
	"github.com/d-ashesss/mah-moneh/internal/users"
	"github.com/gofrs/uuid"
)

type Service struct {
	db Store
}

func NewService(db Store) *Service {
	return &Service{db: db}
}

func (s *Service) AddIncome(ctx context.Context, u *users.User, month string, currency string, amt float64, desc string) (*Transaction, error) {
	panic("implemente me")
}

func (s *Service) AddTransfer(ctx context.Context, u *users.User, month string, currency string, amt float64, desc string) (*Transaction, error) {
	panic("implemente me")
}

func (s *Service) AddExpense(ctx context.Context, u *users.User, month string, currency string, amt float64, desc string) (*Transaction, error) {
	panic("implemente me")
}

func (s *Service) DeleteTransaction(ctx context.Context, tx *Transaction) error {
	panic("implemente me")
}

func (s *Service) GetTransaction(ctx context.Context, uuid uuid.UUID) (*Transaction, error) {
	panic("implemente me")
}

func (s *Service) GetUserTransactions(ctx context.Context, u *users.User, month string) (TransactionCollection, error) {
	panic("implemente me")
}
