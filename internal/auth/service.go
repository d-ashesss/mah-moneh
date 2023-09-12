package auth

import (
	"context"
	"fmt"
	"github.com/d-ashesss/mah-moneh/internal/users"
	"github.com/gofrs/uuid"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"log"
)

type UsersService interface {
	GetUser(ctx context.Context, UUID uuid.UUID) (*users.User, error)
}

type Service struct {
	keySet jwk.Set
	users  UsersService
}

func NewService(cfg *Config, usersSrv UsersService) *Service {
	set := jwk.NewSet()
	if cfg.OpenIDConfigurationUrl != "" {
		log.Printf("[AUTH] Fetching keys from %s", cfg.OpenIDConfigurationUrl)
		var err error
		set, err = jwk.Fetch(context.Background(), cfg.OpenIDConfigurationUrl)
		if err != nil {
			panic(err)
		}
	}
	return &Service{
		keySet: set,
		users:  usersSrv,
	}
}

func (s *Service) AddKey(k jwk.Key) error {
	return s.keySet.AddKey(k)
}

func (s *Service) AuthenticateUser(ctx context.Context, token string) (*users.User, error) {
	t, err := jwt.Parse([]byte(token), jwt.WithKeySet(s.keySet))
	if err != nil {
		return nil, err
	}
	UUID, err := uuid.FromString(t.Subject())
	if err != nil {
		return nil, fmt.Errorf("invalid auth payload: %s", err)
	}
	return s.users.GetUser(ctx, UUID)
}
