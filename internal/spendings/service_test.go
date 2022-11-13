package spendings_test

import (
	"context"
	"github.com/d-ashesss/mah-moneh/internal/accounts"
	"github.com/d-ashesss/mah-moneh/internal/capital"
	"github.com/d-ashesss/mah-moneh/internal/categories"
	mocks "github.com/d-ashesss/mah-moneh/internal/mocks/spendings"
	"github.com/d-ashesss/mah-moneh/internal/spendings"
	"github.com/d-ashesss/mah-moneh/internal/transactions"
	"github.com/d-ashesss/mah-moneh/internal/users"
	"github.com/stretchr/testify/suite"
	"testing"
)

type CapitalServiceTestSuite struct {
	suite.Suite
	capital      *mocks.CapitalService
	transactions *mocks.TransactionsService
	categories   *mocks.CategoryService
	srv          *spendings.Service
}

func (ts *CapitalServiceTestSuite) SetupTest() {
	ts.capital = mocks.NewCapitalService(ts.T())
	ts.transactions = mocks.NewTransactionsService(ts.T())
	ts.categories = mocks.NewCategoryService(ts.T())
	ts.srv = spendings.NewService(ts.capital, ts.transactions, ts.categories)
}

func (ts *CapitalServiceTestSuite) TestGetMonthSpendings() {
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
	catIncome := &categories.Category{}
	catEmpty := &categories.Category{}
	catSomething := &categories.Category{}
	ts.categories.On("GetUserCategories", ctx, u).Return([]*categories.Category{catIncome, catEmpty, catSomething}, nil)
	txs := transactions.TransactionCollection{
		&transactions.Transaction{Amount: -8, Currency: "usd"},
		&transactions.Transaction{Amount: 5, Currency: "usd", Category: catIncome},
		&transactions.Transaction{Amount: -3, Currency: "eur"},
		&transactions.Transaction{Amount: -2, Currency: "eur", Category: catSomething},
		&transactions.Transaction{Amount: -4, Currency: "eth"},
		&transactions.Transaction{Amount: 2, Currency: "btc", Category: catIncome},
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

func TestCapitalService(t *testing.T) {
	suite.Run(t, new(CapitalServiceTestSuite))
}
