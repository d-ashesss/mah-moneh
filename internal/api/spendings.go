package api

import (
	"context"
	"github.com/d-ashesss/mah-moneh/internal/spendings"
	"github.com/d-ashesss/mah-moneh/internal/users"
)

func (s *Service) GetUserMonthSpendings(ctx context.Context, u *users.User, month string) (spendings.Spendings, error) {
	return s.spendings.GetMonthSpendings(ctx, u, month)
}
