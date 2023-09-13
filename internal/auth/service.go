package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/d-ashesss/mah-moneh/internal/users"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"io"
	"log"
	"net/http"
	"time"
)

type UsersService interface {
	GetUser(ctx context.Context, ID string) (*users.User, error)
}

type openidConfiguration struct {
	JwksUri string `json:"jwks_uri"`
}

type Service struct {
	keySet jwk.Set
	users  UsersService
}

func NewService(cfg *Config, usersSrv UsersService) *Service {
	set := jwk.NewSet()
	if cfg.OpenIDConfigurationUrl != "" {
		openidCfg, err := fetchOpenIDConf(cfg.OpenIDConfigurationUrl)
		if err != nil {
			panic(err)
		}

		if openidCfg.JwksUri != "" {
			log.Printf("[AUTH] Fetching keys from %s", openidCfg.JwksUri)
			set, err = jwk.Fetch(context.Background(), openidCfg.JwksUri)
			if err != nil {
				panic(err)
			}
		}
	}
	return &Service{
		keySet: set,
		users:  usersSrv,
	}
}

func fetchOpenIDConf(url string) (*openidConfiguration, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	buf, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}
	var cfg openidConfiguration
	if err := json.Unmarshal(buf, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (s *Service) AddKey(k jwk.Key) error {
	return s.keySet.AddKey(k)
}

func (s *Service) AuthenticateUser(ctx context.Context, token string) (*users.User, error) {
	t, err := jwt.Parse([]byte(token), jwt.WithKeySet(s.keySet))
	if err != nil {
		return nil, err
	}
	if t.Subject() == "" {
		return nil, fmt.Errorf("invalid auth payload")
	}
	return s.users.GetUser(ctx, t.Subject())
}
