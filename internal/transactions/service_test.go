package transactions_test

import (
	"context"
	mocks "github.com/d-ashesss/mah-moneh/internal/mocks/transactions"
	"github.com/d-ashesss/mah-moneh/internal/transactions"
	"github.com/d-ashesss/mah-moneh/internal/users"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type TransactionsServiceTestSuite struct {
	suite.Suite
	store *mocks.Store
	srv   *transactions.Service
}

func (ts *TransactionsServiceTestSuite) SetupTest() {
	ts.store = mocks.NewStore(ts.T())
	ts.srv = transactions.NewService(ts.store)
}

func (ts *TransactionsServiceTestSuite) TestCreateTransaction() {
	ctx := context.Background()
	u := &users.User{}
	ts.store.On("SaveTransaction", ctx, mock.AnythingOfType("*transactions.Transaction")).
		Return(nil)
	tx, err := ts.srv.CreateTransaction(ctx, u, "2010-10", "usd", 10, "test income transaction", nil)
	ts.Require().NoError(err, "Failed to add income transaction.")
	ts.Require().NotNil(tx)
}

func (ts *TransactionsServiceTestSuite) TestDeleteTransaction() {
	ctx := context.Background()
	tx := &transactions.Transaction{}
	ts.store.On("DeleteTransaction", ctx, tx).Return(nil)
	err := ts.srv.DeleteTransaction(ctx, tx)
	ts.Require().NoError(err, "Failed to delete the transaction.")
}

func (ts *TransactionsServiceTestSuite) TestGetTransaction() {
	ctx := context.Background()
	UUID, _ := uuid.NewV4()
	protoTx := &transactions.Transaction{}
	ts.store.On("GetTransaction", ctx, UUID).Return(protoTx, nil)
	tx, err := ts.srv.GetTransaction(ctx, UUID)
	ts.Require().NoError(err, "Failed to get the transaction.")
	ts.Equal(protoTx, tx)
}

func (ts *TransactionsServiceTestSuite) TestGetUserTransactions() {
	ctx := context.Background()
	u := &users.User{}
	ts.store.On("GetUserTransactions", ctx, u, "2010-10").Return(transactions.TransactionCollection{}, nil)
	txs, err := ts.srv.GetUserTransactions(ctx, u, "2010-10")
	ts.Require().NoError(err, "Failed to get user transactions.")
	ts.Require().NotNil(txs, "Invalid transactions response.")
}

func TestTransactionService(t *testing.T) {
	suite.Run(t, new(TransactionsServiceTestSuite))
}
