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
	User      *users.User
	YearMonth string
	Type      TransactionType
	Currency  string
	Amount    float64
}

type TransactionCollection []*Transaction
