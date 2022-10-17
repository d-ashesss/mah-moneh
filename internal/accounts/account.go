package accounts

import (
	"github.com/d-ashesss/mah-moneh/internal/datastore"
	"github.com/d-ashesss/mah-moneh/internal/users"
)

// Account represents an account entity.
type Account struct {
	datastore.Model
	User *users.User `gorm:"embedded;embeddedPrefix:user_"`
	Name string      `gorm:"notNull"`
}

// AccountCollection represents a collection of account entities.
type AccountCollection []*Account
