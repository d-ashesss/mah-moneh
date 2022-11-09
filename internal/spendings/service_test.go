package spendings_test

import (
	"context"
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
	prevCap := &capital.Capital{Amounts: map[string]float64{
		"usd": 13,
		"eur": 20,
		"eth": 4,
	}}
	currentCap := &capital.Capital{Amounts: map[string]float64{
		"usd": 10,
		"eur": 12,
		"btc": 2,
	}}
	ts.capital.On("GetCapital", ctx, u, "2009-12").Return(prevCap, nil)
	ts.capital.On("GetCapital", ctx, u, "2010-01").Return(currentCap, nil)
	txs := transactions.TransactionCollection{
		&transactions.Transaction{Type: transactions.TypeExpense, Amount: -8, Currency: "usd"},
		&transactions.Transaction{Type: transactions.TypeIncome, Amount: 5, Currency: "usd", Tags: []string{"income"}},
		&transactions.Transaction{Type: transactions.TypeExpense, Amount: -3, Currency: "eur"},
		&transactions.Transaction{Type: transactions.TypeExpense, Amount: -2, Currency: "eur", Tags: []string{"some", "thing"}},
		&transactions.Transaction{Type: transactions.TypeTransfer, Amount: -4, Currency: "eth"},
		&transactions.Transaction{Type: transactions.TypeIncome, Amount: 2, Currency: "btc", Tags: []string{"income"}},
	}
	ts.transactions.On("GetUserTransactions", ctx, u, "2010-01").Return(txs, nil)
	catIncome := &categories.Category{Tags: []string{"income"}}
	catNothing := &categories.Category{Tags: []string{"nothing"}}
	catSome := &categories.Category{Tags: []string{"some"}}
	catSome2 := &categories.Category{Tags: []string{"some"}}
	ts.categories.On("GetUserCategories", ctx, u).Return([]*categories.Category{catIncome, catNothing, catSome, catSome2}, nil)
	spending, err := ts.srv.GetMonthSpendings(ctx, u, "2010-01")
	ts.Require().NoError(err, "Failed to get spendings.")

	ts.InDelta(5.0, spending.ByCategory[catIncome]["usd"], 0.001)
	ts.InDelta(0.0, spending.ByCategory[catIncome]["eur"], 0.001)
	ts.InDelta(0.0, spending.ByCategory[catIncome]["eth"], 0.001)
	ts.InDelta(2.0, spending.ByCategory[catIncome]["btc"], 0.001)

	ts.InDelta(0.0, spending.ByCategory[catNothing]["usd"], 0.001)
	ts.InDelta(0.0, spending.ByCategory[catNothing]["eur"], 0.001)
	ts.InDelta(0.0, spending.ByCategory[catNothing]["eth"], 0.001)
	ts.InDelta(0.0, spending.ByCategory[catNothing]["btc"], 0.001)

	ts.InDelta(0.0, spending.ByCategory[catSome]["usd"], 0.001)
	ts.InDelta(-2.0, spending.ByCategory[catSome]["eur"], 0.001)
	ts.InDelta(0.0, spending.ByCategory[catSome]["eth"], 0.001)
	ts.InDelta(0.0, spending.ByCategory[catSome]["btc"], 0.001)

	// in future catSome2 should be equal to catSome
	ts.InDelta(0.0, spending.ByCategory[catSome2]["usd"], 0.001)
	ts.InDelta(0.0, spending.ByCategory[catSome2]["eur"], 0.001)
	ts.InDelta(0.0, spending.ByCategory[catSome2]["eth"], 0.001)
	ts.InDelta(0.0, spending.ByCategory[catSome2]["btc"], 0.001)

	ts.InDelta(-8.0, spending.Uncategorized["usd"], 0.001)
	ts.InDelta(-3.0, spending.Uncategorized["eur"], 0.001)
	ts.InDelta(-4.0, spending.Uncategorized["eth"], 0.001)
	ts.InDelta(0.0, spending.Uncategorized["btc"], 0.001)

	ts.InDelta(0.0, spending.Unaccounted["usd"], 0.001)
	ts.InDelta(-3.0, spending.Unaccounted["eur"], 0.001)
	ts.InDelta(0.0, spending.Unaccounted["eth"], 0.001)
	ts.InDelta(0.0, spending.Unaccounted["btc"], 0.001)
}

func TestCapitalService(t *testing.T) {
	suite.Run(t, new(CapitalServiceTestSuite))
}
