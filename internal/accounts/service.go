package accounts

import (
	"context"
	"github.com/d-ashesss/mah-moneh/internal/users"
	"github.com/gofrs/uuid"
	"time"
)

const FmtYearMonth = "2006-01"

// Service is a service responsible for managing accounts.
type Service struct {
	db AccountStore
}

// NewService initializes a new accounts service.
func NewService(db AccountStore) *Service {
	return &Service{db: db}
}

func (s *Service) CreateAccount(ctx context.Context, u *users.User, name string) (*Account, error) {
	acc := NewAccount(u, name)
	if err := s.db.CreateAccount(ctx, acc); err != nil {
		return nil, err
	}
	return acc, nil
}

func (s *Service) UpdateAccount(ctx context.Context, acc *Account) error {
	return s.db.UpdateAccount(ctx, acc)
}

func (s *Service) DeleteAccount(ctx context.Context, acc *Account) error {
	return s.db.DeleteAccount(ctx, acc)
}

func (s *Service) GetAccount(ctx context.Context, UUID uuid.UUID) (*Account, error) {
	return s.db.GetAccount(ctx, UUID)
}

func (s *Service) GetUserAccounts(ctx context.Context, u *users.User) (AccountCollection, error) {
	return s.db.GetUserAccounts(ctx, u)
}

func (s *Service) SetAccountAmount(ctx context.Context, acc *Account, month string, currency Currency, amount float64) error {
	return s.db.SetAccountAmount(ctx, acc, month, currency, amount)
}

func (s *Service) SetAccountCurrentAmount(ctx context.Context, acc *Account, currency Currency, amount float64) error {
	month := time.Now().Format(FmtYearMonth)
	return s.db.SetAccountAmount(ctx, acc, month, currency, amount)
}

func (s *Service) GetAccountAmounts(ctx context.Context, acc *Account, month string) (CurrencyAmounts, error) {
	amounts, err := s.db.GetAccountAmounts(ctx, acc, month)
	if err != nil {
		return nil, err
	}
	return amounts.GetCurrencyAmounts(), nil
}

func (s *Service) GetAccountCurrentAmounts(ctx context.Context, acc *Account) (CurrencyAmounts, error) {
	month := time.Now().Format(FmtYearMonth)
	return s.GetAccountAmounts(ctx, acc, month)
}
