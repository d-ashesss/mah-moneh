package api

import (
	"context"
	"github.com/d-ashesss/mah-moneh/internal/accounts"
	"github.com/d-ashesss/mah-moneh/internal/users"
	"github.com/gofrs/uuid"
)

func (s *Service) CreateAccount(ctx context.Context, u *users.User, name string) (*accounts.Account, error) {
	return s.accounts.CreateAccount(ctx, u, name)
}

func (s *Service) GetAccount(ctx context.Context, UUID uuid.UUID) (*accounts.Account, error) {
	return s.accounts.GetAccount(ctx, UUID)
}

func (s *Service) GetUserAccounts(ctx context.Context, u *users.User) (accounts.AccountCollection, error) {
	return s.accounts.GetUserAccounts(ctx, u)
}

func (s *Service) UpdateAccount(ctx context.Context, acc *accounts.Account) error {
	return s.accounts.UpdateAccount(ctx, acc)
}

func (s *Service) DeleteAccount(ctx context.Context, acc *accounts.Account) error {
	return s.accounts.DeleteAccount(ctx, acc)
}

func (s *Service) SetAccountAmount(ctx context.Context, acc *accounts.Account, month string, currency accounts.Currency, amount float64) error {
	return s.accounts.SetAccountAmount(ctx, acc, month, currency, amount)
}

func (s *Service) SetAccountCurrentAmount(ctx context.Context, acc *accounts.Account, currency accounts.Currency, amount float64) error {
	return s.accounts.SetAccountCurrentAmount(ctx, acc, currency, amount)
}

func (s *Service) GetAccountAmount(ctx context.Context, acc *accounts.Account, month string) (accounts.CurrencyAmounts, error) {
	return s.accounts.GetAccountAmounts(ctx, acc, month)
}

func (s *Service) GetAccountCurrentAmount(ctx context.Context, acc *accounts.Account) (accounts.CurrencyAmounts, error) {
	return s.accounts.GetAccountCurrentAmounts(ctx, acc)
}
