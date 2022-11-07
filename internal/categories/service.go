package categories

import "context"

type Service struct {
	db Store
}

func NewService(db Store) *Service {
	return &Service{db: db}
}

func (s *Service) CreateCategory(ctx context.Context, name string, tags []string) (*Category, error) {
	cat := NewCategory(name, tags)
	if err := s.db.SaveCategory(ctx, cat); err != nil {
		return nil, err
	}
	return cat, nil
}

func (s *Service) DeleteCategory(ctx context.Context, cat *Category) error {
	return s.db.DeleteCategory(ctx, cat)
}
