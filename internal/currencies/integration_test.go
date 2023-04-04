//go:build integration

package currencies_test

import (
	"github.com/d-ashesss/mah-moneh/internal/currencies"
	"github.com/d-ashesss/mah-moneh/internal/datastore"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"testing"
)

type CurrenciesIntegrationTestSuite struct {
	suite.Suite
	db  *gorm.DB
	srv *currencies.Service
}

func (ts *CurrenciesIntegrationTestSuite) SetupSuite() {
	dbCfg, err := datastore.NewConfig()
	if err != nil {
		ts.T().Fatalf("Invalid database config: %s", err)
	}
	dbCfg.TablePrefix = "crc_test_"
	db, err := datastore.Open(dbCfg)
	if err != nil {
		ts.T().Fatalf("Failed to connect to the DB: %s", err)
	}

	ts.db = db.Session(&gorm.Session{NewDB: true})
	store := currencies.NewGormStore(db.Session(&gorm.Session{NewDB: true}))
	ts.srv = currencies.NewService(store)

	err = db.Migrator().AutoMigrate(&currencies.Rate{})
	if err != nil {
		ts.T().Fatalf("Failed to migrate required tables: %s", err)
	}
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
