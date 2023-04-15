//go:build integration

package accounts_test

import (
	"context"
	"github.com/d-ashesss/mah-moneh/internal/accounts"
	"github.com/d-ashesss/mah-moneh/internal/datastore"
	"github.com/d-ashesss/mah-moneh/internal/users"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"testing"
	"time"
)

type AccountsIntegrationTestSuite struct {
	suite.Suite
	db  *gorm.DB
	srv *accounts.Service
}

func (ts *AccountsIntegrationTestSuite) SetupSuite() {
	dbCfg, err := datastore.NewConfig()
	if err != nil {
		ts.T().Fatalf("Invalid database config: %s", err)
	}
	dbCfg.TablePrefix = "acc_test_"
	db, err := datastore.Open(dbCfg)
	if err != nil {
		ts.T().Fatalf("Failed to connect to the DB: %s", err)
	}

	ts.db = db.Session(&gorm.Session{NewDB: true})
	store := accounts.NewGormStore(db.Session(&gorm.Session{NewDB: true}))
	ts.srv = accounts.NewService(store)

	err = db.Migrator().AutoMigrate(&accounts.Amount{}, &accounts.Account{})
	if err != nil {
		ts.T().Fatalf("Failed to migrate required tables: %s", err)
	}
}

func (ts *AccountsIntegrationTestSuite) TestCreateAccount() {
	u := ts.createTestingUser()

	acc, err := ts.srv.CreateAccount(context.Background(), u, "test-create-account")
	ts.Require().NoError(err, "Failed to create account.")

	foundAcc := &accounts.Account{}
	err = ts.db.First(foundAcc, "user_uuid = ? AND name = ?", u.UUID, acc.Name).Error
	ts.Require().NoError(err, "Failed to find the created account.")
	ts.Equal(acc.UUID, foundAcc.UUID)
}

func (ts *AccountsIntegrationTestSuite) TestUpdateAccount() {
	ts.Run("Exists", func() {
		u := ts.createTestingUser()
		acc := ts.createTestingAccount(u, "test-create-account")

		acc.Name = "test-update-account"
		err := ts.srv.UpdateAccount(context.Background(), acc)
		ts.Require().NoError(err, "Failed to update account.")

		foundAcc := &accounts.Account{}
		err = ts.db.First(foundAcc, "user_uuid = ? AND name = ?", u.UUID, acc.Name).Error
		ts.Require().NoError(err, "Failed to find the updated account.")
		ts.Equal(acc.UUID, foundAcc.UUID)
	})

	ts.Run("NotExists", func() {
		u := ts.createTestingUser()
		UUID, _ := uuid.NewV4()
		acc := &accounts.Account{Model: datastore.Model{UUID: UUID}, User: u, Name: "test-update-account"}
		err := ts.srv.UpdateAccount(context.Background(), acc)
		ts.Require().NoError(err, "Failed to update account.")

		foundAcc := &accounts.Account{}
		err = ts.db.First(foundAcc, "user_uuid = ? AND name = ?", u.UUID, acc.Name).Error
		ts.Require().ErrorIs(err, gorm.ErrRecordNotFound, "Non existing account was saved during the update.")
	})
}

func (ts *AccountsIntegrationTestSuite) TestDeleteAccount() {
	u := ts.createTestingUser()
	protoAcc := ts.createTestingAccount(u, "test-delete-account")

	err := ts.srv.DeleteAccount(context.Background(), protoAcc)
	ts.Require().NoError(err, "Failed to delete account.")

	foundAcc := &accounts.Account{}
	err = ts.db.First(foundAcc, "uuid = ?", protoAcc.UUID).Error
	ts.Require().ErrorIs(err, gorm.ErrRecordNotFound, "Deleted record was found.")

}

func (ts *AccountsIntegrationTestSuite) TestGetAccount() {
	u := ts.createTestingUser()
	protoAcc := ts.createTestingAccount(u, "test-get-account")

	foundAcc, err := ts.srv.GetAccount(context.Background(), protoAcc.UUID)
	ts.Require().NoError(err, "Failed to get the account.")
	ts.Equal(protoAcc.UUID, foundAcc.UUID)
}

func (ts *AccountsIntegrationTestSuite) TestGetUserAccounts() {
	u1 := ts.createTestingUser()
	ts.createTestingAccount(u1, "test-get-user-accounts-1-1")
	ts.createTestingAccount(u1, "test-get-user-accounts-1-2")
	u2 := ts.createTestingUser()
	ts.createTestingAccount(u2, "test-get-user-accounts-2-1")

	accs, err := ts.srv.GetUserAccounts(context.Background(), u1)
	ts.Require().NoError(err, "Failed to get accounts.")
	ts.Len(accs, 2, "Invalid set of found accounts.")
}

func (ts *AccountsIntegrationTestSuite) TestSetAccountAmount() {
	u := ts.createTestingUser()
	acc := ts.createTestingAccount(u, "test-set-account-amount")

	err := ts.srv.SetAccountCurrentAmount(context.Background(), acc, "usd", 10.99)
	ts.Require().NoError(err, "Failed to set USD amount on the account.")

	amount := &accounts.Amount{}
	err = ts.db.First(amount, "account_uuid = ? AND currency_code = ?", acc.UUID, "usd").Error
	ts.Require().NoError(err, "Failed to get USD amount.")
	ts.InDelta(10.99, amount.Amount, 0.001, "Invalid amount on account.")

	err = ts.srv.SetAccountCurrentAmount(context.Background(), acc, "usd", 12)
	ts.Require().NoError(err, "Failed to change USD amount on the account.")

	err = ts.srv.SetAccountCurrentAmount(context.Background(), acc, "eur", 21)
	ts.Require().NoError(err, "Failed to set EUR amount on the account.")

	month := time.Now().Format(accounts.FmtYearMonth)

	amount = &accounts.Amount{}
	err = ts.db.First(amount, "account_uuid = ? AND year_month = ? AND currency_code = ?", acc.UUID, month, "usd").Error
	ts.Require().NoError(err, "Failed to get updated USD amount.")
	ts.InDelta(12.0, amount.Amount, 0.001, "Invalid amount on account.")

	amount = &accounts.Amount{}
	err = ts.db.First(amount, "account_uuid = ? AND year_month = ? AND currency_code = ?", acc.UUID, month, "eur").Error
	ts.Require().NoError(err, "Failed to get EUR amount.")
	ts.InDelta(21.0, amount.Amount, 0.001, "Invalid amount on account.")
}

func (ts *AccountsIntegrationTestSuite) TestGetAccountAmounts() {
	u := ts.createTestingUser()
	acc := ts.createTestingAccount(u, "test-get-account-amounts")
	var (
		amount  *accounts.Amount
		amounts accounts.CurrencyAmounts
		err     error
	)

	amounts, err = ts.srv.GetAccountAmounts(context.Background(), acc, "2010-11")
	ts.Require().NoError(err, "Failed to get amounts on the account.")
	ts.Require().Len(amounts, 0, "Invalid set of amounts returned.")

	amount = &accounts.Amount{Account: acc, YearMonth: "2010-10", CurrencyCode: "usd", Amount: 10.0}
	err = ts.db.Save(amount).Error
	ts.Require().NoError(err, "Failed to set amount on the account.")

	amount = &accounts.Amount{Account: acc, YearMonth: "2010-10", CurrencyCode: "eur", Amount: 15.0}
	err = ts.db.Save(amount).Error
	ts.Require().NoError(err, "Failed to set amount on the account.")

	amount = &accounts.Amount{Account: acc, YearMonth: "2010-08", CurrencyCode: "usd", Amount: 20.0}
	err = ts.db.Save(amount).Error
	ts.Require().NoError(err, "Failed to set amount on the account.")

	amount = &accounts.Amount{Account: acc, YearMonth: "2010-06", CurrencyCode: "usd", Amount: 30.0}
	err = ts.db.Save(amount).Error
	ts.Require().NoError(err, "Failed to set amount on the account.")

	amount = &accounts.Amount{Account: acc, YearMonth: "2010-05", CurrencyCode: "eur", Amount: 35.0}
	err = ts.db.Save(amount).Error
	ts.Require().NoError(err, "Failed to set amount on the account.")

	amounts, err = ts.srv.GetAccountAmounts(context.Background(), acc, "2010-11")
	ts.Require().NoError(err, "Failed to get amounts on the account.")
	ts.Require().Len(amounts, 2, "Invalid set of amounts returned.")
	ts.InDelta(10.0, amounts["usd"], 0.001, "Invalid amount on account.")
	ts.InDelta(15.0, amounts["eur"], 0.001, "Invalid amount on account.")

	amounts, err = ts.srv.GetAccountAmounts(context.Background(), acc, "2010-10")
	ts.Require().NoError(err, "Failed to get amounts on the account.")
	ts.Require().Len(amounts, 2, "Invalid set of amounts returned.")
	ts.InDelta(10.0, amounts["usd"], 0.001, "Invalid amount on account.")
	ts.InDelta(15.0, amounts["eur"], 0.001, "Invalid amount on account.")

	amounts, err = ts.srv.GetAccountAmounts(context.Background(), acc, "2010-09")
	ts.Require().NoError(err, "Failed to get amounts on the account.")
	ts.Require().Len(amounts, 2, "Invalid set of amounts returned.")
	ts.InDelta(20.0, amounts["usd"], 0.001, "Invalid amount on account.")
	ts.InDelta(35.0, amounts["eur"], 0.001, "Invalid amount on account.")

	amounts, err = ts.srv.GetAccountAmounts(context.Background(), acc, "2010-05")
	ts.Require().NoError(err, "Failed to get amounts on the account.")
	ts.Require().Len(amounts, 1, "Invalid set of amounts returned.")
	ts.InDelta(0.0, amounts["usd"], 0.001, "Invalid amount on account.")
	ts.InDelta(35.0, amounts["eur"], 0.001, "Invalid amount on account.")
}

func (ts *AccountsIntegrationTestSuite) TestGetAccountCurrentAmounts() {
	u := ts.createTestingUser()
	acc := ts.createTestingAccount(u, "test-set-account-amount")
	var (
		amount  *accounts.Amount
		amounts accounts.CurrencyAmounts
		err     error
	)

	amounts, err = ts.srv.GetAccountCurrentAmounts(context.Background(), acc)
	ts.Require().NoError(err, "Failed to get amounts on the account.")
	ts.Len(amounts, 0, "Invalid set of amounts returned.")

	amount = &accounts.Amount{Account: acc, YearMonth: "2000-01", CurrencyCode: "usd", Amount: 5.0}
	err = ts.db.Save(amount).Error
	ts.Require().NoError(err, "Failed to set amount on the account.")

	amounts, err = ts.srv.GetAccountCurrentAmounts(context.Background(), acc)
	ts.Require().NoError(err, "Failed to get amounts on the account.")
	ts.Require().Len(amounts, 1, "Invalid set of amounts returned.")
	ts.InDelta(5.0, amounts["usd"], 0.001, "Invalid amount on account.")

	month := time.Now().Format(accounts.FmtYearMonth)

	amount = &accounts.Amount{Account: acc, YearMonth: month, CurrencyCode: "usd", Amount: 11.5}
	err = ts.db.Save(amount).Error
	ts.Require().NoError(err, "Failed to set amount on the account.")

	amount = &accounts.Amount{Account: acc, YearMonth: month, CurrencyCode: "eur", Amount: 27.3}
	err = ts.db.Save(amount).Error
	ts.Require().NoError(err, "Failed to set amount on the account.")

	amounts, err = ts.srv.GetAccountCurrentAmounts(context.Background(), acc)
	ts.Require().NoError(err, "Failed to get amounts on the account.")
	ts.Require().Len(amounts, 2, "Invalid set of amounts returned.")
	ts.InDelta(11.5, amounts["usd"], 0.001, "Invalid amount on account.")
	ts.InDelta(27.3, amounts["eur"], 0.001, "Invalid amount on account.")
}

func (ts *AccountsIntegrationTestSuite) createTestingUser() *users.User {
	ts.T().Helper()
	UUID, _ := uuid.NewV4()
	return &users.User{UUID: UUID}
}

func (ts *AccountsIntegrationTestSuite) createTestingAccount(u *users.User, name string) *accounts.Account {
	ts.T().Helper()
	acc := accounts.NewAccount(u, name)
	err := ts.db.Create(acc).Error
	ts.Require().NoError(err, "Failed to create testing account.")
	return acc
}

func TestAccountsIntegration(t *testing.T) {
	suite.Run(t, new(AccountsIntegrationTestSuite))
}
