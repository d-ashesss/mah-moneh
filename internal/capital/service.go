package capital

import (
	"context"
	"github.com/d-ashesss/mah-moneh/internal/accounts"
	"github.com/d-ashesss/mah-moneh/internal/users"
)

type AccountsService interface {
	GetUserAccounts(ctx context.Context, u *users.User) (accounts.AccountCollection, error)
	GetAccountAmounts(ctx context.Context, acc *accounts.Account) (accounts.AmountCollection, error)
}

type Service struct {
	accounts AccountsService
}

func NewService(as AccountsService) *Service {
	return &Service{accounts: as}
}

func (s *Service) GetCapital(ctx context.Context, u *users.User) (*Capital, error) {
	accs, err := s.accounts.GetUserAccounts(ctx, u)
	if err != nil {
		return nil, err
	}
	c := &Capital{}
	for _, acc := range accs {
		amounts, err := s.accounts.GetAccountAmounts(ctx, acc)
		if err != nil {
			return nil, err
		}
		for _, a := range amounts {
			c.Amount += a.Amount
		}
	}
	return c, nil
}
