package transactions_test

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type TransactionsIntegrationTestSuite struct {
	suite.Suite
}

func (ts *TransactionsIntegrationTestSuite) TestAddIncome() {
	panic("implemente me")
}

func (ts *TransactionsIntegrationTestSuite) TestAddTransfer() {
	panic("implemente me")
}

func (ts *TransactionsIntegrationTestSuite) TestAddExpense() {
	panic("implemente me")
}

func (ts *TransactionsIntegrationTestSuite) TestDeleteTransaction() {
	panic("implemente me")
}

func (ts *TransactionsIntegrationTestSuite) TestGetTransaction() {
	panic("implemente me")
}

func (ts *TransactionsIntegrationTestSuite) TestGetUserTransactions() {
	panic("implemente me")
}

func TestTransactionsIntegration(t *testing.T) {
	suite.Run(t, new(TransactionsIntegrationTestSuite))
}
