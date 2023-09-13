package users

import (
	"context"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) GetUser(_ context.Context, ID string) (*User, error) {
	return &User{ID: ID}, nil
}
