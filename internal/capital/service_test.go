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
		&accounts.Amount{Amount: 10},
		&accounts.Amount{Amount: 12},
	}
	ts.accounts.On("GetUserAccounts", ctx, u).Return(accs, nil)
	ts.accounts.On("GetAccountAmounts", ctx, acc).Return(amounts, nil)
	c, err := ts.srv.GetCapital(ctx, u)
	ts.Require().NoError(err, "Failed to get capital.")
	ts.Equal(22., c.Amount)
}

func TestCapitalSuite(t *testing.T) {
	suite.Run(t, new(CapitalTestSuite))
}
