package transactions_test

import (
	"context"
	"fmt"
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
	ts.db = db
	store := transactions.NewGormStore(db)
	ts.srv = transactions.NewService(store)
}

func (ts *TransactionsIntegrationTestSuite) SetupTest() {
	_ = ts.db.Migrator().DropTable(&transactions.Transaction{})
	_ = ts.db.AutoMigrate(&transactions.Transaction{})
}

func (ts *TransactionsIntegrationTestSuite) TestAddIncome() {
	u := ts.createTestingUser()
	tx, err := ts.srv.AddIncome(context.Background(), u, "2010-10", "usd", 10, "test add income")
	ts.Require().NoError(err, "Failed to create income transaction.")
	ts.Require().NotNil(tx, "Failed to create income transaction.")

	foundTx := &transactions.Transaction{}
	err = ts.db.First(foundTx, "uuid = ?", tx.UUID).Error
	ts.Require().NoError(err, "Failed to find created transaction")
	ts.Equal(tx.UUID, foundTx.UUID)
	ts.InDelta(tx.Amount, foundTx.Amount, 0.001)
	ts.Equal(tx.Description, foundTx.Description)
	ts.Equal(transactions.TypeIncome, foundTx.Type)
}

func (ts *TransactionsIntegrationTestSuite) TestAddTransfer() {
	u := ts.createTestingUser()
	tx, err := ts.srv.AddTransfer(context.Background(), u, "2010-10", "usd", 10, "test add transfer")
	ts.Require().NoError(err, "Failed to create transfer transaction.")
	ts.Require().NotNil(tx, "Failed to create transfer transaction.")

	foundTx := &transactions.Transaction{}
	err = ts.db.First(foundTx, "uuid = ?", tx.UUID).Error
	ts.Require().NoError(err, "Failed to find created transaction")
	ts.Equal(tx.UUID, foundTx.UUID)
	ts.InDelta(tx.Amount, foundTx.Amount, 0.001)
	ts.Equal(tx.Description, foundTx.Description)
	ts.Equal(transactions.TypeTransfer, foundTx.Type)
}

func (ts *TransactionsIntegrationTestSuite) TestAddExpense() {
	u := ts.createTestingUser()
	tx, err := ts.srv.AddExpense(context.Background(), u, "2010-10", "usd", 10, "test add expense")
	ts.Require().NoError(err, "Failed to create expense transaction.")
	ts.Require().NotNil(tx, "Failed to create expense transaction.")

	foundTx := &transactions.Transaction{}
	err = ts.db.First(foundTx, "uuid = ?", tx.UUID).Error
	ts.Require().NoError(err, "Failed to find created transaction")
	ts.Equal(tx.UUID, foundTx.UUID)
	ts.InDelta(tx.Amount, foundTx.Amount, 0.001)
	ts.Equal(tx.Description, foundTx.Description)
	ts.Equal(transactions.TypeExpense, foundTx.Type)
}

func (ts *TransactionsIntegrationTestSuite) TestDeleteTransaction() {
	u := ts.createTestingUser()
	tx := transactions.NewIncomeTransaction(u, "2010-10", "usd", 10, "test delete tx")
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
	tx := transactions.NewIncomeTransaction(u, "2010-10", "usd", 10, "test get tx")
	err := ts.db.Save(tx).Error
	ts.Require().NoError(err, "Failed to save the transaction.")

	foundTx, err := ts.srv.GetTransaction(context.Background(), tx.UUID)
	ts.Require().NoError(err, "Failed to find created transaction")
	ts.Equal(tx.UUID, foundTx.UUID)
	ts.InDelta(tx.Amount, foundTx.Amount, 0.001)
	ts.Equal(tx.Description, foundTx.Description)
	ts.Equal(tx.Type, foundTx.Type)
}

func (ts *TransactionsIntegrationTestSuite) TestGetUserTransactions() {
	u1 := ts.createTestingUser()
	u2 := ts.createTestingUser()
	var (
		tx  *transactions.Transaction
		txs transactions.TransactionCollection
		err error
	)

	tx = transactions.NewIncomeTransaction(u1, "2010-11", "usd", 10, "test tx")
	err = ts.db.Save(tx).Error
	ts.Require().NoError(err, "Failed to save the transaction.")

	tx = transactions.NewIncomeTransaction(u1, "2010-10", "usd", 10, "test tx")
	err = ts.db.Save(tx).Error
	ts.Require().NoError(err, "Failed to save the transaction.")

	tx = transactions.NewIncomeTransaction(u1, "2010-10", "usd", 10, "test tx")
	err = ts.db.Save(tx).Error
	ts.Require().NoError(err, "Failed to save the transaction.")

	tx = transactions.NewIncomeTransaction(u1, "2010-09", "usd", 10, "test tx")
	err = ts.db.Save(tx).Error
	ts.Require().NoError(err, "Failed to save the transaction.")

	tx = transactions.NewIncomeTransaction(u2, "2010-10", "usd", 10, "test tx")
	err = ts.db.Save(tx).Error
	ts.Require().NoError(err, "Failed to save the transaction.")

	txs, err = ts.srv.GetUserTransactions(context.Background(), u1, "2010-10")
	ts.Require().NoError(err, "Failed to get user's transactions.")
	ts.Len(txs, 2)

	txs, err = ts.srv.GetUserTransactions(context.Background(), u1, "2010-09")
	ts.Require().NoError(err, "Failed to get user's transactions.")
	ts.Len(txs, 1)
}

func (ts *TransactionsIntegrationTestSuite) createTestingUser() *users.User {
	ts.T().Helper()
	UUID, _ := uuid.NewV4()
	return &users.User{UUID: UUID}
}

func TestTransactionsIntegration(t *testing.T) {
	suite.Run(t, new(TransactionsIntegrationTestSuite))
}
