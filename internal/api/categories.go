package api

import (
	"context"
	"github.com/d-ashesss/mah-moneh/internal/categories"
	"github.com/d-ashesss/mah-moneh/internal/users"
	"github.com/gofrs/uuid"
)

func (s *Service) CreateCategory(ctx context.Context, u *users.User, name string) (*categories.Category, error) {
	return s.categories.CreateCategory(ctx, u, name)
}

func (s *Service) GetCategory(ctx context.Context, UUID uuid.UUID) (*categories.Category, error) {
	return s.categories.GetCategory(ctx, UUID)
}

func (s *Service) GetUserCategories(ctx context.Context, u *users.User) ([]*categories.Category, error) {
	return s.categories.GetUserCategories(ctx, u)
}

func (s *Service) DeleteCategory(ctx context.Context, cat *categories.Category) error {
	return s.categories.DeleteCategory(ctx, cat)
}
