package capital

import (
	"context"
	"github.com/d-ashesss/mah-moneh/internal/accounts"
	"github.com/d-ashesss/mah-moneh/internal/users"
)

type AccountsService interface {
	GetUserAccounts(ctx context.Context, u *users.User) (accounts.AccountCollection, error)
	GetAccountAmounts(ctx context.Context, acc *accounts.Account, month string) (accounts.AmountCollection, error)
}

// Service is a service responsible for calculating the capital.
type Service struct {
	accounts AccountsService
}

// NewService initializes capital service.
func NewService(as AccountsService) *Service {
	return &Service{accounts: as}
}

// GetCapital calculates capital for the specified month.
func (s *Service) GetCapital(ctx context.Context, u *users.User, month string) (*Capital, error) {
	accs, err := s.accounts.GetUserAccounts(ctx, u)
	if err != nil {
		return nil, err
	}
	c := New()
	for _, acc := range accs {
		amounts, err := s.accounts.GetAccountAmounts(ctx, acc, month)
		if err != nil {
			return nil, err
		}
		for currency, a := range amounts {
			c.Amounts[currency] += a.Amount
		}
	}
	return c, nil
}
