package auth

import "github.com/joeshaw/envdecode"

type Config struct {
	OpenIDConfigurationUrl string `env:"AUTH_OPENID_CONFIGURATION_URL"`
}

func NewConfig() *Config {
	cfg := Config{}
	_ = envdecode.Decode(&cfg)
	return &cfg
}
