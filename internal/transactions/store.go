package transactions

import (
	"context"
	"errors"
	"github.com/d-ashesss/mah-moneh/internal/datastore"
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
	return s.db.WithContext(ctx).Save(tx).Error
}

func (s *gormStore) DeleteTransaction(ctx context.Context, tx *Transaction) error {
	return s.db.WithContext(ctx).Delete(tx).Error
}

func (s *gormStore) GetTransaction(ctx context.Context, uuid uuid.UUID) (*Transaction, error) {
	tx := &Transaction{}
	err := s.db.WithContext(ctx).Preload("Category").First(tx, "uuid = ?", uuid).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, datastore.ErrRecordNotFound
	}
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (s *gormStore) GetUserTransactions(ctx context.Context, u *users.User, month string) (TransactionCollection, error) {
	txs := make(TransactionCollection, 0)
	err := s.db.WithContext(ctx).Preload("Category").Where("user_uuid = ?", u.UUID).Where("year_month = ?", month).Find(&txs).Error
	if err != nil {
		return nil, err
	}
	return txs, nil
}
