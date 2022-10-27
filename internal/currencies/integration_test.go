//go:build integration

package currencies_test

import (
	"fmt"
	"github.com/d-ashesss/mah-moneh/internal/currencies"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
	"testing"
)

type CurrenciesIntegrationTestSuite struct {
	suite.Suite
	db  *gorm.DB
	srv *currencies.Service
}

func (ts *CurrenciesIntegrationTestSuite) SetupSuite() {
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
	store := currencies.NewGormStore(db)
	ts.srv = currencies.NewService(store)
}

func (ts *CurrenciesIntegrationTestSuite) SetupTest() {
	_ = ts.db.Migrator().DropTable(&currencies.Rate{})
	_ = ts.db.AutoMigrate(&currencies.Rate{})
}

func (ts *CurrenciesIntegrationTestSuite) TestSetRate() {
	var (
		err  error
		rate *currencies.Rate
	)
	err = ts.srv.SetRate("usd", "eur", "2010-10", 0.9)
	ts.Require().NoError(err, "Failed to set the rate.")

	rate = &currencies.Rate{}
	err = ts.db.Where("base = ? AND target = ? AND year_month = ?", "usd", "eur", "2010-10").First(rate).Error
	ts.Require().NoError(err, "Failed to get the rate.")
	ts.InDelta(0.9, rate.Rate, 0.001)

	err = ts.srv.SetRate("usd", "eur", "2010-10", 1.1)
	ts.Require().NoError(err, "Failed to update the rate.")

	rate = &currencies.Rate{}
	err = ts.db.Where("base = ? AND target = ? AND year_month = ?", "usd", "eur", "2010-10").First(rate).Error
	ts.Require().NoError(err, "Failed to get the updated rate.")
	ts.InDelta(1.1, rate.Rate, 0.001)
}

func (ts *CurrenciesIntegrationTestSuite) TestGetRate() {
	ts.createRate("usd", "eur", "2010-10", 1.1)
	ts.createRate("usd", "eur", "2010-08", 1.0)
	var rate float64

	rate = ts.srv.GetRate("usd", "eur", "2010-07")
	ts.InDelta(1.0, rate, 0.001)

	rate = ts.srv.GetRate("usd", "eur", "2010-08")
	ts.InDelta(1.0, rate, 0.001)

	rate = ts.srv.GetRate("usd", "eur", "2010-09")
	ts.InDelta(1.0, rate, 0.001)

	rate = ts.srv.GetRate("usd", "eur", "2010-11")
	ts.InDelta(1.1, rate, 0.001)

	rate = ts.srv.GetRate("usd", "eth", "2010-11")
	ts.InDelta(0.0, rate, 0.001)
}

func (ts *CurrenciesIntegrationTestSuite) createRate(base, target, month string, rate float64) {
	ts.T().Helper()
	r := &currencies.Rate{
		Base:      base,
		Target:    target,
		YearMonth: month,
		Rate:      rate,
	}
	err := ts.db.Save(r).Error
	ts.Require().NoError(err, "Failed to create testing rate record.")
}

func TestCurrenciesIntegration(t *testing.T) {
	suite.Run(t, new(CurrenciesIntegrationTestSuite))
}
