package capital_test

import (
	"context"
	"github.com/d-ashesss/mah-moneh/internal/accounts"
	"github.com/d-ashesss/mah-moneh/internal/capital"
	mocks "github.com/d-ashesss/mah-moneh/internal/mocks/capital"
	"github.com/d-ashesss/mah-moneh/internal/users"
	"github.com/stretchr/testify/suite"
	"testing"
)

type CapitalTestSuite struct {
	suite.Suite
	accounts *mocks.AccountsService
	srv      *capital.Service
}

func (ts *CapitalTestSuite) SetupTest() {
	ts.accounts = mocks.NewAccountsService(ts.T())
	ts.srv = capital.NewService(ts.accounts)
}

func (ts *CapitalTestSuite) TestGetCapital() {
	ctx := context.Background()
	u := &users.User{}
	acc := &accounts.Account{}
	accs := accounts.AccountCollection{acc}
	amounts := accounts.AmountCollection{
		"usd": &accounts.Amount{CurrencyCode: "usd", Amount: 10},
		"eur": &accounts.Amount{CurrencyCode: "eur", Amount: 12},
	}
	ts.accounts.On("GetUserAccounts", ctx, u).Return(accs, nil)
	ts.accounts.On("GetAccountAmounts", ctx, acc, "2010-10").Return(amounts, nil)
	c, err := ts.srv.GetCapital(ctx, u, "2010-10")
	ts.Require().NoError(err, "Failed to get capital.")
	ts.Equal(10., c.Amounts["usd"])
	ts.Equal(12., c.Amounts["eur"])
}

func TestCapitalSuite(t *testing.T) {
	suite.Run(t, new(CapitalTestSuite))
}
