package capital

import "github.com/d-ashesss/mah-moneh/internal/accounts"

// Capital is a summary of all funds.
type Capital struct {
	Amounts accounts.CurrencyAmounts
}

// New initializes new capital.
func New() *Capital {
	amounts := accounts.NewCurrencyAmounts()
	return &Capital{Amounts: amounts}
}

// Diff gets the difference from provided capital.
func (c *Capital) Diff(from *Capital) accounts.CurrencyAmounts {
	return c.Amounts.Diff(from.Amounts)
}
