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
	tx := NewIncomeTransaction(u, month, currency, amt, desc)
	if err := s.db.SaveTransaction(ctx, tx); err != nil {
		return nil, err
	}
	return tx, nil
}

func (s *Service) AddTransfer(ctx context.Context, u *users.User, month string, currency string, amt float64, desc string) (*Transaction, error) {
	tx := NewTransferTransaction(u, month, currency, amt, desc)
	if err := s.db.SaveTransaction(ctx, tx); err != nil {
		return nil, err
	}
	return tx, nil
}

func (s *Service) AddExpense(ctx context.Context, u *users.User, month string, currency string, amt float64, desc string) (*Transaction, error) {
	tx := NewExpenseTransaction(u, month, currency, amt, desc)
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
