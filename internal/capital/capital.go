package capital

// Capital is a summary of all funds.
type Capital struct {
	Amounts map[string]float64
}

func New() *Capital {
	amounts := make(map[string]float64)
	return &Capital{Amounts: amounts}
}
