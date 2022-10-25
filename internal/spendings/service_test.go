package spendings_test

import (
	"context"
	"github.com/d-ashesss/mah-moneh/internal/capital"
	mocks "github.com/d-ashesss/mah-moneh/internal/mocks/spendings"
	"github.com/d-ashesss/mah-moneh/internal/spendings"
	"github.com/d-ashesss/mah-moneh/internal/users"
	"github.com/stretchr/testify/suite"
	"testing"
)

type CapitalTestSuite struct {
	suite.Suite
	capital *mocks.CapitalService
	srv     *spendings.Service
}

func (ts *CapitalTestSuite) SetupTest() {
	ts.capital = mocks.NewCapitalService(ts.T())
	ts.srv = spendings.NewService(ts.capital)
}

func (ts *CapitalTestSuite) TestGetMonthSpendings() {
	ctx := context.Background()
	u := &users.User{}
	prevCap := &capital.Capital{Amounts: map[string]float64{
		"usd": 13,
		"eur": 20,
		"eth": 4,
	}}
	currentCap := &capital.Capital{Amounts: map[string]float64{
		"usd": 10,
		"eur": 12,
		"btc": 2,
	}}
	ts.capital.On("GetCapital", ctx, u, "2009-12").Return(prevCap, nil)
	ts.capital.On("GetCapital", ctx, u, "2010-01").Return(currentCap, nil)
	spending, err := ts.srv.GetMonthSpendings(ctx, u, "2010-01")
	ts.Require().NoError(err, "Failed to get spendings.")
	ts.Equal(-3., spending.Amounts["usd"])
	ts.Equal(-8., spending.Amounts["eur"])
	ts.Equal(-4., spending.Amounts["eth"])
	ts.Equal(2., spending.Amounts["btc"])
}

func TestCapitalSuite(t *testing.T) {
	suite.Run(t, new(CapitalTestSuite))
}
