package accounts

import (
	"context"
	"github.com/d-ashesss/mah-moneh/internal/users"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

// AccountStore is an interface for accounts DB API.
type AccountStore interface {
	// CreateAccount saves account entity into the DB.
	CreateAccount(ctx context.Context, acc *Account) error
	// UpdateAccount updates existing account entity.
	UpdateAccount(ctx context.Context, acc *Account) error
	// DeleteAccount deletes account entity from the DB.
	DeleteAccount(ctx context.Context, acc *Account) error
	// GetAccount retrieves accounts by its UUID.
	GetAccount(ctx context.Context, UUID uuid.UUID) (*Account, error)
	// GetUserAccounts retrieves all user accounts.
	GetUserAccounts(ctx context.Context, u *users.User) (AccountCollection, error)
}

// gormStore is GORM implementation of AccountStore.
type gormStore struct {
	db *gorm.DB
}

// NewGormStore initializes GORM implementation of AccountStore.
func NewGormStore(db *gorm.DB) AccountStore {
	return &gormStore{db: db}
}

func (s *gormStore) CreateAccount(ctx context.Context, acc *Account) error {
	return s.db.WithContext(ctx).Create(acc).Error
}

func (s *gormStore) UpdateAccount(ctx context.Context, acc *Account) error {
	return s.db.WithContext(ctx).Where("uuid = ?", acc.UUID).Updates(acc).Error
}

func (s *gormStore) DeleteAccount(ctx context.Context, acc *Account) error {
	return s.db.WithContext(ctx).Delete(acc).Error
}

func (s *gormStore) GetAccount(ctx context.Context, UUID uuid.UUID) (*Account, error) {
	acc := &Account{}
	if err := s.db.WithContext(ctx).Where("uuid = ?", UUID).First(acc).Error; err != nil {
		return nil, err
	}
	return acc, nil
}

func (s *gormStore) GetUserAccounts(ctx context.Context, u *users.User) (AccountCollection, error) {
	accs := make(AccountCollection, 0)
	if err := s.db.WithContext(ctx).Find(&accs, "user_uuid = ?", u.UUID).Error; err != nil {
		return nil, err
	}
	return accs, nil
}
