package capital

// Capital is a summary of all funds.
type Capital struct {
	Amounts map[string]float64
}

func New() *Capital {
	amounts := make(map[string]float64)
	return &Capital{Amounts: amounts}
}

func (c *Capital) Substract(c2 *Capital) map[string]float64 {
	amounts := make(map[string]float64)
	for currency, amount := range c.Amounts {
		amounts[currency] = amount
	}
	for currency, amount := range c2.Amounts {
		amounts[currency] -= amount
	}
	return amounts
}
