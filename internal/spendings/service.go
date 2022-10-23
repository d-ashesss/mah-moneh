package spendings

import (
	"context"
	"github.com/d-ashesss/mah-moneh/internal/accounts"
	"github.com/d-ashesss/mah-moneh/internal/capital"
	"github.com/d-ashesss/mah-moneh/internal/users"
	"time"
)

type CapitalService interface {
	GetCapital(ctx context.Context, u *users.User, month string) (*capital.Capital, error)
}

type Service struct {
	capital CapitalService
}

func NewService(cs CapitalService) *Service {
	return &Service{capital: cs}
}

func (s *Service) GetSpendings(ctx context.Context, u *users.User, month string) (*Spending, error) {
	currentCapital, err := s.capital.GetCapital(ctx, u, month)
	if err != nil {
		return nil, err
	}
	prevMonth, err := getPrevMonth(month)
	if err != nil {
		return nil, err
	}
	prevCapital, err := s.capital.GetCapital(ctx, u, prevMonth)
	if err != nil {
		return nil, err
	}
	sp := &Spending{
		Amounts: currentCapital.Substract(prevCapital),
	}
	return sp, nil
}

func getPrevMonth(month string) (string, error) {
	d, err := time.Parse(accounts.FmtYearMonth, month)
	if err != nil {
		return "", err
	}
	return d.AddDate(0, -1, 0).Format(accounts.FmtYearMonth), nil
}
