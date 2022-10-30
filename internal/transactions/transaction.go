package transactions

import (
	"github.com/d-ashesss/mah-moneh/internal/datastore"
	"github.com/d-ashesss/mah-moneh/internal/users"
)

type TransactionType string

const (
	TypeIncome   TransactionType = "income"
	TypeTransfer TransactionType = "transfer"
	TypeExpense  TransactionType = "expense"
)

type Transaction struct {
	datastore.Model
	User        *users.User `gorm:"embedded;embeddedPrefix:user_;notNull"`
	YearMonth   string
	Type        TransactionType
	Currency    string
	Amount      float64
	Description string
}

func NewIncomeTransaction(u *users.User, month string, currency string, amt float64, desc string) *Transaction {
	return &Transaction{
		User:        u,
		YearMonth:   month,
		Type:        TypeIncome,
		Currency:    currency,
		Amount:      amt,
		Description: desc,
	}
}

func NewTransferTransaction(u *users.User, month string, currency string, amt float64, desc string) *Transaction {
	return &Transaction{
		User:        u,
		YearMonth:   month,
		Type:        TypeTransfer,
		Currency:    currency,
		Amount:      amt,
		Description: desc,
	}
}

func NewExpenseTransaction(u *users.User, month string, currency string, amt float64, desc string) *Transaction {
	return &Transaction{
		User:        u,
		YearMonth:   month,
		Type:        TypeExpense,
		Currency:    currency,
		Amount:      amt,
		Description: desc,
	}
}

type TransactionCollection []*Transaction
