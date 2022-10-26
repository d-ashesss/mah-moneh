package converter

type CurrencyService interface {
	GetRate(from, target string) float64
}

// Service represents converter service.
type Service struct {
	currencies CurrencyService
}

// NewService initializes new converter service.
func NewService(cs CurrencyService) *Service {
	return &Service{currencies: cs}
}

// GetTotal calculates total amount in specified currency.
func (s *Service) GetTotal(amounts map[string]float64, targetCurrency string) float64 {
	var total float64
	for currency, amount := range amounts {
		rate := s.currencies.GetRate(currency, targetCurrency)
		total += amount * rate
	}
	return total
}
