package model

import (
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
	"time"
)

// Model defines fields common for most models.
type Model struct {
	UUID      uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

// User represents a user.
type User struct {
	UUID uuid.UUID `gorm:"type:uuid;notNull"`
}

// Account represents an account.
type Account struct {
	Model
	User *User  `gorm:"embedded;embeddedPrefix:user_"`
	Name string `gorm:"notNull"`
}

// NewAccount initializes new account.
func NewAccount(u *User, name string) *Account {
	return &Account{Model: Model{UUID: uuid.Nil}, User: u, Name: name}
}
