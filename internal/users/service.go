package users

import (
	"context"
	"github.com/gofrs/uuid"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) GetUser(_ context.Context, UUID uuid.UUID) (*User, error) {
	return &User{UUID: UUID}, nil
}
