package api

import (
	"github.com/d-ashesss/mah-moneh/internal/accounts"
)

type Service struct {
	accountsSrv *accounts.Service
}

func NewService(accountsSrv *accounts.Service) *Service {
	return &Service{accountsSrv: accountsSrv}
}
