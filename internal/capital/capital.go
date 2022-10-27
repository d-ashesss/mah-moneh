package capital

// Capital is a summary of all funds.
type Capital struct {
	Amounts map[string]float64
}

// New initializes new capital.
func New() *Capital {
	amounts := make(map[string]float64)
	return &Capital{Amounts: amounts}
}

// Subtract subtracts one capital from another.
func (c *Capital) Subtract(c2 *Capital) map[string]float64 {
	amounts := make(map[string]float64)
	for currency, amount := range c.Amounts {
		amounts[currency] = amount
	}
	for currency, amount := range c2.Amounts {
		amounts[currency] -= amount
	}
	return amounts
}
