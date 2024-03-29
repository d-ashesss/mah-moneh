package currencies_test

import (
	"errors"
	"github.com/d-ashesss/mah-moneh/internal/accounts"
	"github.com/d-ashesss/mah-moneh/internal/currencies"
	mocks "github.com/d-ashesss/mah-moneh/internal/mocks/currencies"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type CurrenciesServiceTestSuite struct {
	suite.Suite
	store *mocks.Store
	srv   *currencies.Service
}

func (ts *CurrenciesServiceTestSuite) SetupTest() {
	ts.store = mocks.NewStore(ts.T())
	ts.srv = currencies.NewService(ts.store)
}

func (ts *CurrenciesServiceTestSuite) TestSetRate() {
	ts.store.On("SetRate", accounts.Currency("usd"), accounts.Currency("eur"), "2010-10", 10.0).
		Return(nil).Once()
	err := ts.srv.SetRate("usd", "eur", "2010-10", 10)
	ts.Require().NoError(err, "Failed to set the rate.")
}

func (ts *CurrenciesServiceTestSuite) TestGetRate() {
	eurRate := &currencies.Rate{Rate: 1.1}
	ts.store.On("GetRate", accounts.Currency("usd"), accounts.Currency("eur"), "2010-10").
		Return(eurRate, nil)
	ts.store.On("GetRate", mock.AnythingOfType("accounts.Currency"), mock.AnythingOfType("accounts.Currency"), mock.AnythingOfType("string")).
		Return(nil, errors.New("not found")).Maybe()

	eur := ts.srv.GetRate("usd", "eur", "2010-10")
	ts.InDelta(1.1, eur, 0.001, "Got invalid rate.")

	eth := ts.srv.GetRate("usd", "eth", "2010-10")
	ts.InDelta(0., eth, 0.001, "Got invalid rate.")
}

func TestCurrenciesService(t *testing.T) {
	suite.Run(t, new(CurrenciesServiceTestSuite))
}
