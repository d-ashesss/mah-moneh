package categories

import "context"

type Service struct {
	db Store
}

func NewService(db Store) *Service {
	return &Service{db: db}
}

func (s *Service) CreateCategory(ctx context.Context, name string, tags []string) (*Category, error) {
	panic("implement CreateCategory")
}

func (s *Service) DeleteCategory(ctx context.Context, category *Category) error {
	panic("implement DeleteCategory")
}
