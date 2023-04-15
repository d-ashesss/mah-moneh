package currencies

import "github.com/d-ashesss/mah-moneh/internal/accounts"

// Service represents currencies service.
type Service struct {
	db Store
}

// NewService initializes new currencies service.
func NewService(db Store) *Service {
	return &Service{db: db}
}

// SetRate sets the conversion rate for requested currencies in specified month.
func (s *Service) SetRate(base, target accounts.Currency, month string, rate float64) error {
	return s.db.SetRate(base, target, month, rate)
}

// GetRate provides the conversion rate for requested currencies in specified month.
func (s *Service) GetRate(base, target accounts.Currency, month string) float64 {
	r, err := s.db.GetRate(base, target, month)
	if err != nil {
		return 0
	}
	return r.Rate
}
