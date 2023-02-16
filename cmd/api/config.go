package main

import (
	"os"
	"time"
)

type Config struct {
	Port            string
	ShutdownTimeout time.Duration
}

func LoadConfig() *Config {
	c := &Config{
		Port:            "8080",
		ShutdownTimeout: 60 * time.Second,
	}
	if port, ok := os.LookupEnv("PORT"); ok {
		c.Port = port
	}
	return c
}
