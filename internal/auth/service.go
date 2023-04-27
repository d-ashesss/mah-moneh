package auth

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"github.com/d-ashesss/mah-moneh/internal/users"
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type UsersService interface {
	GetUser(ctx context.Context, UUID uuid.UUID) (*users.User, error)
}

type Service struct {
	key   *rsa.PublicKey
	users UsersService
}

func NewService(cfg *Config, usersSrv UsersService) *Service {
	key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(cfg.PublicKey))
	if err != nil {
		panic(err)
	}
	return &Service{
		key:   key,
		users: usersSrv,
	}
}

func (s *Service) AuthenticateUser(ctx context.Context, token string) (*users.User, error) {
	t, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return s.key, nil
	}, jwt.WithLeeway(5*time.Second))
	if err != nil {
		return nil, err
	}
	claims, ok := t.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return nil, errors.New("invalid auth payload")
	}
	UUID, err := uuid.FromString(claims.Subject)
	if err != nil {
		return nil, fmt.Errorf("invalid auth payload: %s", err)
	}
	return s.users.GetUser(ctx, UUID)
}
