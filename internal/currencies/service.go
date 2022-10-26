package currencies

// Service represents currencies service.
type Service struct {
	db Store
}

// NewService initializes new currencies service.
func NewService(db Store) *Service {
	return &Service{db: db}
}

// GetRate provides the conversion rate for requested currencies in specified month.
func (s *Service) GetRate(base, target, month string) float64 {
	r, err := s.db.GetRate(base, target, month)
	if err != nil {
		return 0
	}
	return r.Rate
}
