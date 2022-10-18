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

type ServiceTestSuite struct {
	suite.Suite
	store *mocks.AccountStore
	srv   *accounts.Accounts
}

func (s *ServiceTestSuite) SetupTest() {
	s.store = mocks.NewAccountStore(s.T())
	s.srv = accounts.NewService(s.store)
}

func (s *ServiceTestSuite) TestCreateAccount() {
	ctx := context.Background()
	s.store.On("CreateAccount", ctx, mock.AnythingOfType("*accounts.Account")).
		Return(nil).Once()

	acc := &accounts.Account{}
	err := s.srv.CreateAccount(ctx, acc)
	s.Require().NoError(err, "Failed to create account.")
}

func (s *ServiceTestSuite) TestUpdateAccount() {
	ctx := context.Background()
	s.store.On("UpdateAccount", ctx, mock.AnythingOfType("*accounts.Account")).
		Return(nil).Once()

	acc := &accounts.Account{}
	err := s.srv.UpdateAccount(ctx, acc)
	s.Require().NoError(err, "Failed to update account.")
}

func (s *ServiceTestSuite) TestDeleteAccount() {
	ctx := context.Background()
	s.store.On("DeleteAccount", ctx, mock.AnythingOfType("*accounts.Account")).
		Return(nil).Once()

	acc := &accounts.Account{}
	err := s.srv.DeleteAccount(ctx, acc)
	s.Require().NoError(err, "Failed to delete account.")
}

func (s *ServiceTestSuite) TestGetAccount() {
	ctx := context.Background()
	UUID, _ := uuid.NewV4()
	protoAcc := &accounts.Account{Model: datastore.Model{UUID: UUID}}
	s.store.On("GetAccount", ctx, UUID).
		Return(protoAcc, nil).Once()

	acc, err := s.srv.GetAccount(ctx, UUID)
	s.Require().NoError(err, "Failed to get the account.")
	s.Equal(protoAcc, acc)
}

func (s *ServiceTestSuite) TestGetUserAccounts() {
	ctx := context.Background()
	s.store.On("GetUserAccounts", ctx, mock.AnythingOfType("*users.User")).
		Return(accounts.AccountCollection{}, nil).Once()

	u := &users.User{}
	accs, err := s.srv.GetUserAccounts(ctx, u)
	s.Require().NoError(err, "Failed to get user accounts.")
	s.NotNil(accs)
}

func TestService(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}
