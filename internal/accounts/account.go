package accounts

import (
	"github.com/d-ashesss/mah-moneh/internal/datastore"
	"github.com/d-ashesss/mah-moneh/internal/users"
	"github.com/gofrs/uuid"
)

// Account represents an account entity.
type Account struct {
	datastore.Model
	User *users.User `gorm:"embedded;embeddedPrefix:user_"`
	Name string      `gorm:"notNull"`
}

// NewAccount initializes a new account.
func NewAccount(u *users.User, name string) *Account {
	return &Account{User: u, Name: name}
}

// AccountCollection represents a collection of account entities.
type AccountCollection []*Account

// Amount represents account's amount entity.
type Amount struct {
	AccountUUID  uuid.UUID `gorm:"primaryKey"`
	Account      *Account  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	YearMonth    string    `gorm:"primaryKey;type:varchar(7);notNull"`
	CurrencyCode string    `gorm:"primaryKey;notNull"`
	Amount       float64
}

// AmountCollection represents a collection of account's amounts.
type AmountCollection map[string]*Amount

func NewAmountCollection(amounts []*Amount) AmountCollection {
	c := make(AmountCollection)
	for _, a := range amounts {
		c[a.CurrencyCode] = a
	}
	return c
}
