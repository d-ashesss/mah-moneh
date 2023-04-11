package accounts_test

import (
	"github.com/d-ashesss/mah-moneh/internal/accounts"
	"github.com/stretchr/testify/suite"
	"testing"
)

type AccountTestSuite struct {
	suite.Suite
}

func (ts *AccountTestSuite) TestAmountCollection_GetCurrencyAmounts_Empty() {
	c := accounts.AmountCollection{}
	ts.Len(c, 0)
}

func (ts *AccountTestSuite) TestAmountCollection_GetCurrencyAmounts() {
	c := accounts.AmountCollection{
		&accounts.Amount{CurrencyCode: "usd", Amount: 5},
		&accounts.Amount{CurrencyCode: "eur", Amount: -2.2},
		&accounts.Amount{CurrencyCode: "eur", Amount: 2.2},
	}
	got := c.GetCurrencyAmounts()
	ts.InDelta(5, got["usd"], 0.001)
	ts.InDelta(0, got["eur"], 0.001)
	ts.InDelta(0, got["btc"], 0.001)
}

func (ts *AccountTestSuite) TestCurrencyAmounts_Diff() {
	a1 := accounts.CurrencyAmounts{"usd": 100, "btc": 0.03}
	a2 := accounts.CurrencyAmounts{"usd": 11.2, "eur": 5}
	got := a1.Diff(a2)
	want := accounts.CurrencyAmounts{"usd": 88.8, "eur": -5, "btc": 0.03}
	ts.Equal(want, got)
}

func TestAccount(t *testing.T) {
	suite.Run(t, new(AccountTestSuite))
}
