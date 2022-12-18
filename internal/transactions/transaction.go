package transactions

import (
	"github.com/d-ashesss/mah-moneh/internal/categories"
	"github.com/d-ashesss/mah-moneh/internal/datastore"
	"github.com/d-ashesss/mah-moneh/internal/users"
	"github.com/gofrs/uuid"
)

type Transaction struct {
	datastore.Model
	User         *users.User `gorm:"embedded;embeddedPrefix:user_;notNull;index"`
	YearMonth    string
	Currency     string
	Amount       float64
	Description  string
	CategoryUUID *uuid.UUID           `gorm:"index"`
	Category     *categories.Category `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

func NewTransaction(u *users.User, month string, currency string, amt float64, desc string, cat *categories.Category) *Transaction {
	return &Transaction{
		User:        u,
		YearMonth:   month,
		Currency:    currency,
		Amount:      amt,
		Description: desc,
		Category:    cat,
	}
}

type TransactionCollection []*Transaction
