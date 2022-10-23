package spendings

type Spending struct {
	Amounts map[string]float64
}

func New() *Spending {
	a := make(map[string]float64)
	return &Spending{Amounts: a}
}
