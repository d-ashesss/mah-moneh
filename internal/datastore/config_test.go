package datastore_test

import (
	"github.com/d-ashesss/mah-moneh/internal/datastore"
	"testing"
)

func TestConfig_Dsn(t *testing.T) {
	tests := []struct {
		name string
		cfg  *datastore.Config
		want string
	}{
		{
			name: "default",
			cfg: &datastore.Config{
				Host: "localhost",
				Name: "test_db",
			},
			want: "host=localhost database=test_db",
		},
		{
			name: "port",
			cfg: &datastore.Config{
				Host: "localhost",
				Port: "5432",
				Name: "test_db",
			},
			want: "host=localhost port=5432 database=test_db",
		},
		{
			name: "credentials",
			cfg: &datastore.Config{
				Host:     "localhost",
				User:     "test_user",
				Password: "test_pwd",
				Name:     "test_db",
			},
			want: "host=localhost user=test_user password=test_pwd database=test_db",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cfg.Dsn(); got != tt.want {
				t.Errorf("Dsn() = %q, want %q", got, tt.want)
			}
		})
	}
}
