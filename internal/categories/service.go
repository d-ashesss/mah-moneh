package categories

import (
	"context"
	"github.com/d-ashesss/mah-moneh/internal/users"
	"github.com/gofrs/uuid"
)

type Service struct {
	db Store
}

func NewService(db Store) *Service {
	return &Service{db: db}
}

func (s *Service) CreateCategory(ctx context.Context, u *users.User, name string) (*Category, error) {
	cat := NewCategory(u, name)
	if err := s.db.SaveCategory(ctx, cat); err != nil {
		return nil, err
	}
	return cat, nil
}

func (s *Service) GetCategory(ctx context.Context, UUID uuid.UUID) (*Category, error) {
	return s.db.GetCategory(ctx, UUID)
}

func (s *Service) DeleteCategory(ctx context.Context, cat *Category) error {
	return s.db.DeleteCategory(ctx, cat)
}

func (s *Service) GetUserCategories(ctx context.Context, u *users.User) ([]*Category, error) {
	return s.db.GetUserCategories(ctx, u)
}
