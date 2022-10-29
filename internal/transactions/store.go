package transactions

import (
	"context"
	"github.com/d-ashesss/mah-moneh/internal/users"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type Store interface {
	SaveTransaction(ctx context.Context, tx *Transaction) error
	DeleteTransaction(ctx context.Context, tx *Transaction) error
	GetTransaction(ctx context.Context, uuid uuid.UUID) (*Transaction, error)
	GetUserTransactions(ctx context.Context, u *users.User, month string) (TransactionCollection, error)
}

type gormStore struct {
	db *gorm.DB
}

func NewGormStore(db *gorm.DB) Store {
	return &gormStore{db: db}
}

func (s *gormStore) SaveTransaction(ctx context.Context, tx *Transaction) error {
	panic("implement me")
}

func (s *gormStore) DeleteTransaction(ctx context.Context, tx *Transaction) error {
	panic("implement me")
}

func (s *gormStore) GetTransaction(ctx context.Context, uuid uuid.UUID) (*Transaction, error) {
	panic("implement me")
}

func (s *gormStore) GetUserTransactions(ctx context.Context, u *users.User, month string) (TransactionCollection, error) {
	panic("implement me")
}
