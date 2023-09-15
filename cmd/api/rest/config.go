package rest

import "github.com/joeshaw/envdecode"

type Config struct {
	AllowedOrigins []string `env:"CORS_ALLOWED_ORIGINS"`
}

func NewConfig() *Config {
	cfg := Config{
		AllowedOrigins: []string{},
	}
	_ = envdecode.Decode(&cfg)
	return &cfg
}
