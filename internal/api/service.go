package api

import (
	"github.com/d-ashesss/mah-moneh/internal/accounts"
	"github.com/d-ashesss/mah-moneh/internal/categories"
)

type Service struct {
	accountsSrv   *accounts.Service
	categoriesSrv *categories.Service
}

func NewService(
	accountsSrv *accounts.Service,
	categoriesSrv *categories.Service,
) *Service {
	return &Service{
		accountsSrv:   accountsSrv,
		categoriesSrv: categoriesSrv,
	}
}
