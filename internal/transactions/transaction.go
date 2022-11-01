package transactions

import (
	"github.com/d-ashesss/mah-moneh/internal/datastore"
	"github.com/d-ashesss/mah-moneh/internal/users"
	"github.com/lib/pq"
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
	Tags        pq.StringArray `gorm:"type:text[]"`
}

func NewIncomeTransaction(u *users.User, month string, currency string, amt float64, desc string, tags []string) *Transaction {
	return &Transaction{
		User:        u,
		YearMonth:   month,
		Type:        TypeIncome,
		Currency:    currency,
		Amount:      amt,
		Description: desc,
		Tags:        tags,
	}
}

func NewTransferTransaction(u *users.User, month string, currency string, amt float64, desc string, tags []string) *Transaction {
	return &Transaction{
		User:        u,
		YearMonth:   month,
		Type:        TypeTransfer,
		Currency:    currency,
		Amount:      amt,
		Description: desc,
		Tags:        tags,
	}
}

func NewExpenseTransaction(u *users.User, month string, currency string, amt float64, desc string, tags []string) *Transaction {
	return &Transaction{
		User:        u,
		YearMonth:   month,
		Type:        TypeExpense,
		Currency:    currency,
		Amount:      amt,
		Description: desc,
		Tags:        tags,
	}
}

type TransactionCollection []*Transaction
