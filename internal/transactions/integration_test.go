//go:build integration

package transactions_test

import (
	"context"
	"fmt"
	"github.com/d-ashesss/mah-moneh/internal/categories"
	"github.com/d-ashesss/mah-moneh/internal/transactions"
	"github.com/d-ashesss/mah-moneh/internal/users"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
	"testing"
)

type TransactionsIntegrationTestSuite struct {
	suite.Suite
	db  *gorm.DB
	srv *transactions.Service
}

func (ts *TransactionsIntegrationTestSuite) SetupSuite() {
	dbHost := os.Getenv("POSTGRES_HOST")
	dbPort := os.Getenv("POSTGRES_PORT")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPwd := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	if dbHost == "" {
		ts.T().Skip("No DB configuration provided.")
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s database=%s", dbHost, dbUser, dbPwd, dbName)
	if dbPort != "" {
		dsn = fmt.Sprintf("%s port=%s", dsn, dbPort)
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Discard})
	if err != nil {
		ts.T().Fatalf("Failed to connect to the DB: %s", err)
	}
	ts.db = db.Session(&gorm.Session{NewDB: true})
	store := transactions.NewGormStore(db.Session(&gorm.Session{NewDB: true}))
	ts.srv = transactions.NewService(store)

	err = db.Migrator().AutoMigrate(&transactions.Transaction{})
	if err != nil {
		ts.T().Fatalf("Failed to migrate required tables: %s", err)
	}
}

func (ts *TransactionsIntegrationTestSuite) TestCreateTransaction() {
	u := ts.createTestingUser()
	tx, err := ts.srv.CreateTransaction(context.Background(), u, "2010-10", "usd", 10, "test add income", nil)
	ts.Require().NoError(err, "Failed to create income transaction.")
	ts.Require().NotNil(tx, "Failed to create income transaction.")

	foundTx := &transactions.Transaction{}
	err = ts.db.First(foundTx, "uuid = ?", tx.UUID).Error
	ts.Require().NoError(err, "Failed to find created transaction")
	ts.Equal(tx.UUID, foundTx.UUID)
	ts.InDelta(tx.Amount, foundTx.Amount, 0.001)
	ts.Equal(tx.Description, foundTx.Description)
	ts.Nil(foundTx.CategoryUUID)
	ts.Nil(foundTx.Category)
}

func (ts *TransactionsIntegrationTestSuite) TestCreateTransactionWithCategory() {
	u := ts.createTestingUser()
	cat := categories.NewCategory(u, "test-category")
	err := ts.db.Save(cat).Error
	ts.Require().NoError(err, "Failed to save testing category.")

	tx, err := ts.srv.CreateTransaction(context.Background(), u, "2010-10", "usd", 10, "test add income", cat)
	ts.Require().NoError(err, "Failed to create income transaction.")
	ts.Require().NotNil(tx, "Failed to create income transaction.")

	foundTx := &transactions.Transaction{}
	err = ts.db.Preload("Category").First(foundTx, "uuid = ?", tx.UUID).Error
	ts.Require().NoError(err, "Failed to find created transaction")
	ts.Equal(tx.UUID, foundTx.UUID)
	ts.InDelta(tx.Amount, foundTx.Amount, 0.001)
	ts.Equal(tx.Description, foundTx.Description)
	ts.NotEmpty(foundTx.CategoryUUID)
	ts.NotEmpty(foundTx.Category)
	ts.Equal(cat.UUID.String(), foundTx.CategoryUUID.String())
}

func (ts *TransactionsIntegrationTestSuite) TestDeleteTransaction() {
	u := ts.createTestingUser()
	tx := transactions.NewTransaction(u, "2010-10", "usd", 10, "test delete tx", nil)
	err := ts.db.Save(tx).Error
	ts.Require().NoError(err, "Failed to save the transaction.")

	err = ts.srv.DeleteTransaction(context.Background(), tx)
	ts.Require().NoError(err, "Failed to delete the transaction.")

	foundTx := &transactions.Transaction{}
	err = ts.db.First(foundTx, "uuid = ?", tx.UUID).Error
	ts.Require().ErrorIs(err, gorm.ErrRecordNotFound, "Deleted transaction should not be found.")
}

func (ts *TransactionsIntegrationTestSuite) TestGetTransaction() {
	u := ts.createTestingUser()
	tx := transactions.NewTransaction(u, "2010-10", "usd", 10, "test get tx", nil)
	err := ts.db.Save(tx).Error
	ts.Require().NoError(err, "Failed to save the transaction.")

	foundTx, err := ts.srv.GetTransaction(context.Background(), tx.UUID)
	ts.Require().NoError(err, "Failed to find created transaction")
	ts.Equal(tx.UUID, foundTx.UUID)
	ts.InDelta(tx.Amount, foundTx.Amount, 0.001)
	ts.Equal(tx.Description, foundTx.Description)
	ts.Nil(foundTx.CategoryUUID)
	ts.Nil(foundTx.Category)
}

func (ts *TransactionsIntegrationTestSuite) TestGetTransactionWithCategory() {
	u := ts.createTestingUser()
	cat := categories.NewCategory(u, "test-category")
	err := ts.db.Save(cat).Error
	ts.Require().NoError(err, "Failed to save testing category.")

	tx := transactions.NewTransaction(u, "2010-10", "usd", 10, "test get tx", cat)
	err = ts.db.Save(tx).Error
	ts.Require().NoError(err, "Failed to save the transaction.")

	foundTx, err := ts.srv.GetTransaction(context.Background(), tx.UUID)
	ts.Require().NoError(err, "Failed to find created transaction")
	ts.Equal(tx.UUID, foundTx.UUID)
	ts.InDelta(tx.Amount, foundTx.Amount, 0.001)
	ts.Equal(tx.Description, foundTx.Description)
	ts.NotEmpty(foundTx.CategoryUUID)
	ts.NotEmpty(foundTx.Category)
	ts.Equal(cat.UUID.String(), foundTx.CategoryUUID.String())
}

func (ts *TransactionsIntegrationTestSuite) TestGetUserTransactions() {
	u1 := ts.createTestingUser()
	u2 := ts.createTestingUser()
	var (
		tx  *transactions.Transaction
		txs transactions.TransactionCollection
		err error
	)

	cat := categories.NewCategory(u1, "test-category-u1")
	err = ts.db.Save(cat).Error
	ts.Require().NoError(err, "Failed to save testing category.")

	tx = transactions.NewTransaction(u1, "2010-11", "usd", 10, "test tx", nil)
	err = ts.db.Save(tx).Error
	ts.Require().NoError(err, "Failed to save the transaction.")

	tx = transactions.NewTransaction(u1, "2010-10", "usd", 10, "test tx", cat)
	err = ts.db.Save(tx).Error
	ts.Require().NoError(err, "Failed to save the transaction.")

	tx = transactions.NewTransaction(u1, "2010-10", "usd", 10, "test tx", nil)
	err = ts.db.Save(tx).Error
	ts.Require().NoError(err, "Failed to save the transaction.")

	tx = transactions.NewTransaction(u1, "2010-09", "usd", 10, "test tx", cat)
	err = ts.db.Save(tx).Error
	ts.Require().NoError(err, "Failed to save the transaction.")

	tx = transactions.NewTransaction(u2, "2010-10", "usd", 10, "test tx", nil)
	err = ts.db.Save(tx).Error
	ts.Require().NoError(err, "Failed to save the transaction.")

	txs, err = ts.srv.GetUserTransactions(context.Background(), u1, "2010-10")
	ts.Require().NoError(err, "Failed to get user's transactions.")
	ts.Len(txs, 2)
	ts.NotEmpty(txs[0].Category)
	ts.Nil(txs[1].Category)

	txs, err = ts.srv.GetUserTransactions(context.Background(), u1, "2010-09")
	ts.Require().NoError(err, "Failed to get user's transactions.")
	ts.Len(txs, 1)
	ts.NotEmpty(txs[0].Category)

}

func (ts *TransactionsIntegrationTestSuite) createTestingUser() *users.User {
	ts.T().Helper()
	UUID, _ := uuid.NewV4()
	return &users.User{UUID: UUID}
}

func TestTransactionsIntegration(t *testing.T) {
	suite.Run(t, new(TransactionsIntegrationTestSuite))
}
