package main

import (
	"github.com/joeshaw/envdecode"
	"time"
)

type Config struct {
	Port            string `env:"PORT,default=8080"`
	ShutdownTimeout time.Duration
}

func NewConfig() *Config {
	cfg := Config{
		ShutdownTimeout: 60 * time.Second,
	}
	_ = envdecode.Decode(&cfg)
	return &cfg
}
