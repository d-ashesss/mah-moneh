package converter

import "github.com/d-ashesss/mah-moneh/internal/accounts"

type CurrencyService interface {
	GetRate(base, target accounts.Currency, month string) float64
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
func (s *Service) GetTotal(amounts accounts.CurrencyAmounts, targetCurrency accounts.Currency, month string) float64 {
	var total float64
	for currency, amount := range amounts {
		rate := s.currencies.GetRate(currency, targetCurrency, month)
		total += amount * rate
	}
	return total
}
