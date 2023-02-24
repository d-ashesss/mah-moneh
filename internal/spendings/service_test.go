package spendings_test

import (
	"context"
	"github.com/d-ashesss/mah-moneh/internal/accounts"
	"github.com/d-ashesss/mah-moneh/internal/capital"
	"github.com/d-ashesss/mah-moneh/internal/categories"
	"github.com/d-ashesss/mah-moneh/internal/datastore"
	mocks "github.com/d-ashesss/mah-moneh/internal/mocks/spendings"
	"github.com/d-ashesss/mah-moneh/internal/spendings"
	"github.com/d-ashesss/mah-moneh/internal/transactions"
	"github.com/d-ashesss/mah-moneh/internal/users"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/suite"
	"testing"
)

type SpendingsServiceTestSuite struct {
	suite.Suite
	capital      *mocks.CapitalService
	transactions *mocks.TransactionsService
	categories   *mocks.CategoryService
	srv          *spendings.Service
}

func (ts *SpendingsServiceTestSuite) SetupTest() {
	ts.capital = mocks.NewCapitalService(ts.T())
	ts.transactions = mocks.NewTransactionsService(ts.T())
	ts.categories = mocks.NewCategoryService(ts.T())
	ts.srv = spendings.NewService(ts.capital, ts.transactions, ts.categories)
}

func newCategory(UUID string) *categories.Category {
	return &categories.Category{
		Model: datastore.Model{
			UUID: uuid.FromStringOrNil(UUID),
		},
	}
}

func (ts *SpendingsServiceTestSuite) TestGetMonthSpendings() {
	ctx := context.Background()
	u := &users.User{}
	prevCap := &capital.Capital{Amounts: accounts.CurrencyAmounts{
		"usd": 13,
		"eur": 20,
		"eth": 4,
	}}
	currentCap := &capital.Capital{Amounts: accounts.CurrencyAmounts{
		"usd": 10,
		"eur": 12,
		"btc": 2,
	}}
	ts.capital.On("GetCapital", ctx, u, "2009-12").Return(prevCap, nil)
	ts.capital.On("GetCapital", ctx, u, "2010-01").Return(currentCap, nil)
	catIncomeUUID := "7070e309-af27-445a-9b15-3f9db12a5377"
	catIncome := newCategory(catIncomeUUID)
	catEmptyUUID := "01fbfdc3-c399-4cdf-bbf8-37b3422e6466"
	catEmpty := newCategory(catEmptyUUID)
	catSomethingUUID := "8d4357a0-d20b-410b-a520-cf7d575c402f"
	catSomething := newCategory(catSomethingUUID)
	ts.categories.On("GetUserCategories", ctx, u).Return([]*categories.Category{newCategory(catIncomeUUID), newCategory(catEmptyUUID), newCategory(catSomethingUUID)}, nil)
	txs := transactions.TransactionCollection{
		&transactions.Transaction{Amount: -8, Currency: "usd"},
		&transactions.Transaction{Amount: 5, Currency: "usd", Category: newCategory(catIncomeUUID)},
		&transactions.Transaction{Amount: -3, Currency: "eur"},
		&transactions.Transaction{Amount: -2, Currency: "eur", Category: newCategory(catSomethingUUID)},
		&transactions.Transaction{Amount: -4, Currency: "eth"},
		&transactions.Transaction{Amount: 2, Currency: "btc", Category: newCategory(catIncomeUUID)},
	}
	ts.transactions.On("GetUserTransactions", ctx, u, "2010-01").Return(txs, nil)
	spending, err := ts.srv.GetMonthSpendings(ctx, u, "2010-01")
	ts.Require().NoError(err, "Failed to get spendings.")

	ts.InDelta(5.0, spending.GetAmount(catIncome, "usd"), 0.001)
	ts.InDelta(0.0, spending.GetAmount(catIncome, "eur"), 0.001)
	ts.InDelta(0.0, spending.GetAmount(catIncome, "eth"), 0.001)
	ts.InDelta(2.0, spending.GetAmount(catIncome, "btc"), 0.001)

	ts.InDelta(0.0, spending.GetAmount(catEmpty, "usd"), 0.001)
	ts.InDelta(0.0, spending.GetAmount(catEmpty, "eur"), 0.001)
	ts.InDelta(0.0, spending.GetAmount(catEmpty, "eth"), 0.001)
	ts.InDelta(0.0, spending.GetAmount(catEmpty, "btc"), 0.001)

	ts.InDelta(0.0, spending.GetAmount(catSomething, "usd"), 0.001)
	ts.InDelta(-2.0, spending.GetAmount(catSomething, "eur"), 0.001)
	ts.InDelta(0.0, spending.GetAmount(catSomething, "eth"), 0.001)
	ts.InDelta(0.0, spending.GetAmount(catSomething, "btc"), 0.001)

	ts.InDelta(-8.0, spending.GetAmount(spendings.Uncategorized, "usd"), 0.001)
	ts.InDelta(-3.0, spending.GetAmount(spendings.Uncategorized, "eur"), 0.001)
	ts.InDelta(-4.0, spending.GetAmount(spendings.Uncategorized, "eth"), 0.001)
	ts.InDelta(0.0, spending.GetAmount(spendings.Uncategorized, "btc"), 0.001)

	ts.InDelta(0.0, spending.GetAmount(spendings.Unaccounted, "usd"), 0.001)
	ts.InDelta(-3.0, spending.GetAmount(spendings.Unaccounted, "eur"), 0.001)
	ts.InDelta(0.0, spending.GetAmount(spendings.Unaccounted, "eth"), 0.001)
	ts.InDelta(0.0, spending.GetAmount(spendings.Unaccounted, "btc"), 0.001)
}

func TestSpendingsService(t *testing.T) {
	suite.Run(t, new(SpendingsServiceTestSuite))
}
