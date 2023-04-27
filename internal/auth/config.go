package auth

import "github.com/joeshaw/envdecode"

type Config struct {
	PublicKey string `env:"PUBLIC_KEY"`
}

func NewConfig() *Config {
	cfg := Config{}
	_ = envdecode.Decode(&cfg)
	return &cfg
}
