//go:build integration

package accounts_test

import (
	"context"
	"fmt"
	"github.com/d-ashesss/mah-moneh/internal/accounts"
	"github.com/d-ashesss/mah-moneh/internal/datastore"
	"github.com/d-ashesss/mah-moneh/internal/users"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
	"testing"
	"time"
)

type AccountsTestSuite struct {
	suite.Suite
	db  *gorm.DB
	srv *accounts.Accounts
}

func (s *AccountsTestSuite) SetupSuite() {
	dbHost := os.Getenv("POSTGRES_HOST")
	dbPort := os.Getenv("POSTGRES_PORT")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPwd := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	if dbHost == "" {
		s.T().Skip("No DB configuration provided.")
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s database=%s", dbHost, dbUser, dbPwd, dbName)
	if dbPort != "" {
		dsn = fmt.Sprintf("%s port=%s", dsn, dbPort)
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Discard})
	if err != nil {
		s.T().Fatalf("Failed to connect to the DB: %s", err)
	}
	s.db = db
	store := accounts.NewGormStore(db)
	s.srv = accounts.NewService(store)
}

func (s *AccountsTestSuite) SetupTest() {
	_ = s.db.Migrator().DropTable(&accounts.Amount{}, &accounts.Account{})
	_ = s.db.AutoMigrate(&accounts.Account{}, &accounts.Amount{})
}

func (s *AccountsTestSuite) TestCreateAccount() {
	u := s.createTestingUser()
	acc := accounts.NewAccount(u, "test-create-account")

	err := s.srv.CreateAccount(context.Background(), acc)
	s.Require().NoError(err, "Failed to create account.")

	foundAcc := &accounts.Account{}
	err = s.db.First(foundAcc, "user_uuid = ? AND name = ?", u.UUID, acc.Name).Error
	s.Require().NoError(err, "Failed to find the created account.")
	s.Equal(acc.UUID, foundAcc.UUID)
}

func (s *AccountsTestSuite) TestUpdateAccount() {
	s.Run("Exists", func() {
		u := s.createTestingUser()
		acc := s.createTestingAccount(u, "test-create-account")

		acc.Name = "test-update-account"
		err := s.srv.UpdateAccount(context.Background(), acc)
		s.Require().NoError(err, "Failed to update account.")

		foundAcc := &accounts.Account{}
		err = s.db.First(foundAcc, "user_uuid = ? AND name = ?", u.UUID, acc.Name).Error
		s.Require().NoError(err, "Failed to find the updated account.")
		s.Equal(acc.UUID, foundAcc.UUID)
	})

	s.Run("NotExists", func() {
		u := s.createTestingUser()
		UUID, _ := uuid.NewV4()
		acc := &accounts.Account{Model: datastore.Model{UUID: UUID}, User: u, Name: "test-update-account"}
		err := s.srv.UpdateAccount(context.Background(), acc)
		s.Require().NoError(err, "Failed to update account.")

		foundAcc := &accounts.Account{}
		err = s.db.First(foundAcc, "user_uuid = ? AND name = ?", u.UUID, acc.Name).Error
		s.Require().ErrorIs(err, gorm.ErrRecordNotFound, "Non existing account was saved during the update.")
	})
}

func (s *AccountsTestSuite) TestDeleteAccount() {
	u := s.createTestingUser()
	protoAcc := s.createTestingAccount(u, "test-delete-account")

	err := s.srv.DeleteAccount(context.Background(), protoAcc)
	s.Require().NoError(err, "Failed to delete account.")

	foundAcc := &accounts.Account{}
	err = s.db.First(foundAcc, "uuid = ?", protoAcc.UUID).Error
	s.Require().ErrorIs(err, gorm.ErrRecordNotFound, "Deleted record was found.")

}

func (s *AccountsTestSuite) TestGetAccount() {
	u := s.createTestingUser()
	protoAcc := s.createTestingAccount(u, "test-get-account")

	foundAcc, err := s.srv.GetAccount(context.Background(), protoAcc.UUID)
	s.Require().NoError(err, "Failed to get the account.")
	s.Equal(protoAcc.UUID, foundAcc.UUID)
}

func (s *AccountsTestSuite) TestGetUserAccounts() {
	u1 := s.createTestingUser()
	s.createTestingAccount(u1, "test-get-user-accounts-1-1")
	s.createTestingAccount(u1, "test-get-user-accounts-1-2")
	u2 := s.createTestingUser()
	s.createTestingAccount(u2, "test-get-user-accounts-2-1")

	accs, err := s.srv.GetUserAccounts(context.Background(), u1)
	s.Require().NoError(err, "Failed to get accounts.")
	s.Len(accs, 2, "Invalid set of found accounts.")
}

func (s *AccountsTestSuite) TestSetAccountAmount() {
	u := s.createTestingUser()
	acc := s.createTestingAccount(u, "test-set-account-amount")

	err := s.srv.SetAccountCurrentAmount(context.Background(), acc, "usd", 10.99)
	s.Require().NoError(err, "Failed to set USD amount on the account.")

	amount := &accounts.Amount{}
	err = s.db.First(amount, "account_uuid = ? AND currency_code = ?", acc.UUID, "usd").Error
	s.Require().NoError(err, "Failed to get USD amount.")
	s.Equal(10.99, amount.Amount, "Invalid amount on account.")

	err = s.srv.SetAccountCurrentAmount(context.Background(), acc, "usd", 12)
	s.Require().NoError(err, "Failed to change USD amount on the account.")

	err = s.srv.SetAccountCurrentAmount(context.Background(), acc, "eur", 21)
	s.Require().NoError(err, "Failed to set EUR amount on the account.")

	month := time.Now().Format("2006-01")

	amount = &accounts.Amount{}
	err = s.db.First(amount, "account_uuid = ? AND year_month = ? AND currency_code = ?", acc.UUID, month, "usd").Error
	s.Require().NoError(err, "Failed to get updated USD amount.")
	s.Equal(12., amount.Amount, "Invalid amount on account.")

	amount = &accounts.Amount{}
	err = s.db.First(amount, "account_uuid = ? AND year_month = ? AND currency_code = ?", acc.UUID, month, "eur").Error
	s.Require().NoError(err, "Failed to get EUR amount.")
	s.Equal(21., amount.Amount, "Invalid amount on account.")
}

func (s *AccountsTestSuite) TestGetAccountAmounts() {
	u := s.createTestingUser()
	acc := s.createTestingAccount(u, "test-set-account-amount")

	amounts, err := s.srv.GetAccountCurrentAmounts(context.Background(), acc)
	s.Require().NoError(err, "Failed to get amounts on the account.")
	s.Len(amounts, 0, "Invalid set of amounts returned.")

	month := time.Now().Format("2006-01")

	amount := &accounts.Amount{Account: acc, YearMonth: month, CurrencyCode: "usd", Amount: 11.5}
	err = s.db.Save(amount).Error
	s.Require().NoError(err, "Failed to set amount on the account.")

	amount = &accounts.Amount{Account: acc, YearMonth: month, CurrencyCode: "eur", Amount: 27.3}
	err = s.db.Save(amount).Error
	s.Require().NoError(err, "Failed to set amount on the account.")

	amount = &accounts.Amount{Account: acc, YearMonth: "2000-01", CurrencyCode: "usd", Amount: 5.}
	err = s.db.Save(amount).Error
	s.Require().NoError(err, "Failed to set amount on the account.")

	amounts, err = s.srv.GetAccountCurrentAmounts(context.Background(), acc)
	s.Require().NoError(err, "Failed to get amounts on the account.")
	s.Require().Len(amounts, 2, "Invalid set of amounts returned.")
	s.Equal(11.5, amounts["usd"].Amount, "Invalid amount on account.")
	s.Equal(27.3, amounts["eur"].Amount, "Invalid amount on account.")
}

func (s *AccountsTestSuite) createTestingUser() *users.User {
	s.T().Helper()
	UUID, _ := uuid.NewV4()
	return &users.User{UUID: UUID}
}

func (s *AccountsTestSuite) createTestingAccount(u *users.User, name string) *accounts.Account {
	s.T().Helper()
	acc := accounts.NewAccount(u, name)
	err := s.db.Create(acc).Error
	s.Require().NoError(err, "Failed to create testing account.")
	return acc
}

func TestAccountsTestSuite(t *testing.T) {
	suite.Run(t, new(AccountsTestSuite))
}
