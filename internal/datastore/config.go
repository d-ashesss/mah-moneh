package datastore

import (
	"fmt"
	"github.com/joeshaw/envdecode"
)

type Config struct {
	Host     string `env:"DB_HOST,default=localhost"`
	Port     string `env:"DB_PORT,default=5432"`
	User     string `env:"DB_USER,default=postgres"`
	Password string `env:"DB_PASSWORD"`
	Name     string `env:"DB_NAME,default=postgres"`

	Debug bool `env:"DB_DEBUG,default=false"`

	TablePrefix string
}

func NewConfig() (*Config, error) {
	var cfg Config
	if err := envdecode.Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (c *Config) Dsn() string {
	dsn := fmt.Sprintf("host=%s", c.Host)
	if len(c.Port) > 0 {
		dsn = fmt.Sprintf("%s port=%s", dsn, c.Port)
	}
	if len(c.User) > 0 {
		dsn = fmt.Sprintf("%s user=%s", dsn, c.User)
	}
	if len(c.Password) > 0 {
		dsn = fmt.Sprintf("%s password=%s", dsn, c.Password)
	}
	dsn = fmt.Sprintf("%s database=%s", dsn, c.Name)
	return dsn
}
