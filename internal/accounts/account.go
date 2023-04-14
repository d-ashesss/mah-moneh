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
type AmountCollection []*Amount

// GetCurrencyAmounts extracts amounts for each currency in the collection.
func (a AmountCollection) GetCurrencyAmounts() CurrencyAmounts {
	c := NewCurrencyAmounts()
	for _, amt := range a {
		c[amt.CurrencyCode] += amt.Amount
	}
	return c
}

// CurrencyAmounts contains amount per currency.
type CurrencyAmounts map[string]float64

// NewCurrencyAmounts creates new currency amounts instance.
func NewCurrencyAmounts() CurrencyAmounts {
	return make(CurrencyAmounts)
}

// Diff gets the difference from provided amounts.
func (a CurrencyAmounts) Diff(from CurrencyAmounts) CurrencyAmounts {
	diff := NewCurrencyAmounts()
	for currency, amount := range a {
		diff[currency] = amount
	}
	for currency, amount := range from {
		diff[currency] -= amount
	}
	return diff.filterZeroValues()
}

func (a CurrencyAmounts) filterZeroValues() CurrencyAmounts {
	for currency, amount := range a {
		if -0.0001 < amount && amount < 0.0001 {
			delete(a, currency)
		}
	}
	return a
}
