package accounts

import (
	"context"
	"github.com/d-ashesss/mah-moneh/internal/users"
	"github.com/gofrs/uuid"
)

// Accounts is a service responsible for managing accounts.
type Accounts struct {
	db AccountStore
}

// NewService initializes a new accounts service.
func NewService(db AccountStore) *Accounts {
	return &Accounts{db: db}
}

func (s *Accounts) CreateAccount(ctx context.Context, acc *Account) error {
	return s.db.CreateAccount(ctx, acc)
}

func (s *Accounts) UpdateAccount(ctx context.Context, acc *Account) error {
	return s.db.UpdateAccount(ctx, acc)
}

func (s *Accounts) DeleteAccount(ctx context.Context, acc *Account) error {
	return s.db.DeleteAccount(ctx, acc)
}

func (s *Accounts) GetAccount(ctx context.Context, UUID uuid.UUID) (*Account, error) {
	return s.db.GetAccount(ctx, UUID)
}

func (s *Accounts) GetUserAccounts(ctx context.Context, u *users.User) (AccountCollection, error) {
	return s.db.GetUserAccounts(ctx, u)
}

func (s *Accounts) SetAccountAmount(ctx context.Context, acc *Account, currency string, amount float64) error {
	return s.db.SetAccountAmount(ctx, acc, currency, amount)
}

func (s *Accounts) GetAccountAmounts(ctx context.Context, acc *Account) (AmountCollection, error) {
	return s.db.GetAccountAmounts(ctx, acc)
}
