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
