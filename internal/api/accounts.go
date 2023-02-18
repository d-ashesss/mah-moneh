package api

import (
	"context"
	"github.com/d-ashesss/mah-moneh/internal/accounts"
	"github.com/d-ashesss/mah-moneh/internal/users"
	"github.com/gofrs/uuid"
)

func (s *Service) CreateAccount(ctx context.Context, u *users.User, name string) (*accounts.Account, error) {
	acc := &accounts.Account{
		User: u,
		Name: name,
	}
	if err := s.accountsSrv.CreateAccount(ctx, acc); err != nil {
		return nil, err
	}
	return acc, nil
}

func (s *Service) GetAccount(ctx context.Context, UUID uuid.UUID) (*accounts.Account, error) {
	return s.accountsSrv.GetAccount(ctx, UUID)
}

func (s *Service) GetUserAccounts(ctx context.Context, u *users.User) (accounts.AccountCollection, error) {
	return s.accountsSrv.GetUserAccounts(ctx, u)
}

func (s *Service) UpdateAccount(ctx context.Context, acc *accounts.Account) error {
	return s.accountsSrv.UpdateAccount(ctx, acc)
}

func (s *Service) DeleteAccount(ctx context.Context, acc *accounts.Account) error {
	return s.accountsSrv.DeleteAccount(ctx, acc)
}

func (s *Service) SetAccountAmount(ctx context.Context, acc *accounts.Account, month string, currency string, amount float64) error {
	return s.accountsSrv.SetAccountAmount(ctx, acc, month, currency, amount)
}

func (s *Service) SetAccountCurrentAmount(ctx context.Context, acc *accounts.Account, currency string, amount float64) error {
	return s.accountsSrv.SetAccountCurrentAmount(ctx, acc, currency, amount)
}

func (s *Service) GetAccountAmount(ctx context.Context, acc *accounts.Account, month string) (accounts.CurrencyAmounts, error) {
	return s.accountsSrv.GetAccountAmounts(ctx, acc, month)
}

func (s *Service) GetAccountCurrentAmount(ctx context.Context, acc *accounts.Account) (accounts.CurrencyAmounts, error) {
	return s.accountsSrv.GetAccountCurrentAmounts(ctx, acc)
}
