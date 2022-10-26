package currencies

// Service represents currencies service.
type Service struct {
}

// NewService initializes new currencies service.
func NewService() *Service {
	return &Service{}
}

// GetRate provides the conversion rate for requested currencies in specified month.
func (s *Service) GetRate(from, target, month string) float64 {
	return 0
}
