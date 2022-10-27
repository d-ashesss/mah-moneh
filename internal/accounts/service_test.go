package accounts_test

import (
	"context"
	"github.com/d-ashesss/mah-moneh/internal/accounts"
	"github.com/d-ashesss/mah-moneh/internal/datastore"
	mocks "github.com/d-ashesss/mah-moneh/internal/mocks/accounts"
	"github.com/d-ashesss/mah-moneh/internal/users"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type AccountsServiceTestSuite struct {
	suite.Suite
	store *mocks.AccountStore
	srv   *accounts.Service
}

func (ts *AccountsServiceTestSuite) SetupTest() {
	ts.store = mocks.NewAccountStore(ts.T())
	ts.srv = accounts.NewService(ts.store)
}

func (ts *AccountsServiceTestSuite) TestCreateAccount() {
	ctx := context.Background()
	ts.store.On("CreateAccount", ctx, mock.AnythingOfType("*accounts.Account")).
		Return(nil).Once()

	acc := &accounts.Account{}
	err := ts.srv.CreateAccount(ctx, acc)
	ts.Require().NoError(err, "Failed to create account.")
}

func (ts *AccountsServiceTestSuite) TestUpdateAccount() {
	ctx := context.Background()
	ts.store.On("UpdateAccount", ctx, mock.AnythingOfType("*accounts.Account")).
		Return(nil).Once()

	acc := &accounts.Account{}
	err := ts.srv.UpdateAccount(ctx, acc)
	ts.Require().NoError(err, "Failed to update account.")
}

func (ts *AccountsServiceTestSuite) TestDeleteAccount() {
	ctx := context.Background()
	ts.store.On("DeleteAccount", ctx, mock.AnythingOfType("*accounts.Account")).
		Return(nil).Once()

	acc := &accounts.Account{}
	err := ts.srv.DeleteAccount(ctx, acc)
	ts.Require().NoError(err, "Failed to delete account.")
}

func (ts *AccountsServiceTestSuite) TestGetAccount() {
	ctx := context.Background()
	UUID, _ := uuid.NewV4()
	protoAcc := &accounts.Account{Model: datastore.Model{UUID: UUID}}
	ts.store.On("GetAccount", ctx, UUID).
		Return(protoAcc, nil).Once()

	acc, err := ts.srv.GetAccount(ctx, UUID)
	ts.Require().NoError(err, "Failed to get the account.")
	ts.Equal(protoAcc, acc)
}

func (ts *AccountsServiceTestSuite) TestGetUserAccounts() {
	ctx := context.Background()
	ts.store.On("GetUserAccounts", ctx, mock.AnythingOfType("*users.User")).
		Return(accounts.AccountCollection{}, nil).Once()

	u := &users.User{}
	accs, err := ts.srv.GetUserAccounts(ctx, u)
	ts.Require().NoError(err, "Failed to get user accounts.")
	ts.NotNil(accs)
}

func (ts *AccountsServiceTestSuite) TestSetAccountCurrentAmount() {
	ctx := context.Background()
	acc := &accounts.Account{}
	ts.store.On("SetAccountAmount", ctx, acc, mock.AnythingOfType("string"), "usd", 10.).
		Return(nil).Once()

	err := ts.srv.SetAccountCurrentAmount(ctx, acc, "usd", 10)
	ts.Require().NoError(err, "Failed to set amount on the account.")
}

func (ts *AccountsServiceTestSuite) TestGetAccountAmounts() {
	ctx := context.Background()
	acc := &accounts.Account{}
	ts.store.On("GetAccountAmounts", ctx, acc, "2010-01").
		Return(accounts.AmountCollection{}, nil).Once()

	amounts, err := ts.srv.GetAccountAmounts(ctx, acc, "2010-01")
	ts.Require().NoError(err, "Failed to get amounts of the account.")
	ts.NotNil(amounts)
}

func (ts *AccountsServiceTestSuite) TestGetAccountCurrentAmounts() {
	ctx := context.Background()
	acc := &accounts.Account{}
	ts.store.On("GetAccountAmounts", ctx, acc, mock.AnythingOfType("string")).
		Return(accounts.AmountCollection{}, nil).Once()

	amounts, err := ts.srv.GetAccountCurrentAmounts(ctx, acc)
	ts.Require().NoError(err, "Failed to get amounts of the account.")
	ts.NotNil(amounts)
}

func TestAccountsService(t *testing.T) {
	suite.Run(t, new(AccountsServiceTestSuite))
}
