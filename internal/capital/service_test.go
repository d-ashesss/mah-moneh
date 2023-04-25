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

type CapitalServiceTestSuite struct {
	suite.Suite
	accounts *mocks.AccountsService
	srv      *capital.Service
}

func (ts *CapitalServiceTestSuite) SetupTest() {
	ts.accounts = mocks.NewAccountsService(ts.T())
	ts.srv = capital.NewService(ts.accounts)
}

func (ts *CapitalServiceTestSuite) TestGetCapital_SingleAccount() {
	ctx := context.Background()
	u := &users.User{}
	acc := &accounts.Account{}
	accs := accounts.AccountCollection{acc}
	amounts := accounts.CurrencyAmounts{
		"usd": 10,
		"eur": 12,
	}
	ts.accounts.On("GetUserAccounts", ctx, u).Return(accs, nil)
	ts.accounts.On("GetAccountAmounts", ctx, acc, "2010-10").Return(amounts, nil)
	c, err := ts.srv.GetCapital(ctx, u, "2010-10")
	ts.Require().NoError(err, "Failed to get capital.")
	ts.InDelta(10.0, c.Amounts["usd"], 0.001)
	ts.InDelta(12.0, c.Amounts["eur"], 0.001)
}

func (ts *CapitalServiceTestSuite) TestGetCapital_MultipleAccounts() {
	ctx := context.Background()
	u := &users.User{}
	acc1 := &accounts.Account{}
	acc2 := &accounts.Account{}
	accs := accounts.AccountCollection{acc1, acc2}
	amounts1 := accounts.CurrencyAmounts{
		"usd": 15,
	}
	amounts2 := accounts.CurrencyAmounts{
		"usd": 10,
		"eur": 12,
	}
	ts.accounts.On("GetUserAccounts", ctx, u).Return(accs, nil)
	ts.accounts.On("GetAccountAmounts", ctx, acc1, "2010-10").Return(amounts1, nil).Once()
	ts.accounts.On("GetAccountAmounts", ctx, acc2, "2010-10").Return(amounts2, nil).Once()
	c, err := ts.srv.GetCapital(ctx, u, "2010-10")
	ts.Require().NoError(err, "Failed to get capital.")
	ts.InDelta(25.0, c.Amounts["usd"], 0.001)
	ts.InDelta(12.0, c.Amounts["eur"], 0.001)
}

func TestCapitalService(t *testing.T) {
	suite.Run(t, new(CapitalServiceTestSuite))
}
