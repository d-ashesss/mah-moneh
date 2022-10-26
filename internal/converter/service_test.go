package converter_test

import (
	"github.com/d-ashesss/mah-moneh/internal/converter"
	mocks "github.com/d-ashesss/mah-moneh/internal/mocks/converter"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ConverterTestSuite struct {
	suite.Suite
	srv *converter.Service
}

func (ts *ConverterTestSuite) SetupTest() {
	cs := mocks.NewCurrencyService(ts.T())
	cs.On("GetRate", "usd", "usd").Return(1.).Maybe()
	cs.On("GetRate", "usd", "eur").Return(1.1).Maybe()
	cs.On("GetRate", "usd", "btc").Return(0.1).Maybe()
	cs.On("GetRate", "eur", "eur").Return(1.).Maybe()
	cs.On("GetRate", "eur", "usd").Return(0.91).Maybe()
	cs.On("GetRate", "eur", "btc").Return(0.091).Maybe()
	cs.On("GetRate", "btc", "btc").Return(1.).Maybe()
	cs.On("GetRate", "btc", "usd").Return(10.).Maybe()
	cs.On("GetRate", "btc", "eur").Return(10.99).Maybe()
	cs.On("GetRate", mock.Anything, mock.Anything).Return(0.).Maybe()
	ts.srv = converter.NewService(cs)
}

func (ts *ConverterTestSuite) TestGetTotal() {
	amounts := map[string]float64{
		"usd": 100,
		"eur": 100,
		"btc": 5,
		"eth": 15,
	}
	var total float64

	total = ts.srv.GetTotal(amounts, "usd")
	ts.Equal(241., total)

	total = ts.srv.GetTotal(amounts, "eur")
	ts.T().Log(total)
	ts.Equal(264.95, total)

	total = ts.srv.GetTotal(amounts, "btc")
	ts.Equal(24.1, total)
}

func TestConverterTestSuite(t *testing.T) {
	suite.Run(t, new(ConverterTestSuite))
}
