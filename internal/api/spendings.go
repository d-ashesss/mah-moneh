package api

import (
	"context"
	"github.com/d-ashesss/mah-moneh/internal/accounts"
	"github.com/d-ashesss/mah-moneh/internal/spendings"
	"github.com/d-ashesss/mah-moneh/internal/users"
)

type Spendings map[string]accounts.CurrencyAmounts

func (s *Service) GetUserMonthSpendings(ctx context.Context, u *users.User, month string) (Spendings, error) {
	cats, err := s.categories.GetUserCategories(ctx, u)
	if err != nil {
		return nil, err
	}
	spent, err := s.spendings.GetMonthSpendings(ctx, u, month)
	if err != nil {
		return nil, err
	}
	r := make(Spendings)
	for _, cat := range cats {
		r[cat.UUID.String()] = spent.GetAmounts(cat)
	}
	r["uncategorized"] = spent.GetAmounts(spendings.Uncategorized)
	r["unaccounted"] = spent.GetAmounts(spendings.Unaccounted)
	r["total"] = spent.GetAmounts(spendings.Total)
	return r, nil
}
