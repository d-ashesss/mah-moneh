package categories

import (
	"context"
	"github.com/d-ashesss/mah-moneh/internal/users"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type Store interface {
	SaveCategory(ctx context.Context, cat *Category) error
	DeleteCategory(ctx context.Context, cat *Category) error
	GetCategory(ctx context.Context, uuid uuid.UUID) (*Category, error)
	GetUserCategories(ctx context.Context, u *users.User) ([]*Category, error)
}

type gormStore struct {
	db *gorm.DB
}

func NewGormStore(db *gorm.DB) Store {
	return &gormStore{db: db}
}

func (s *gormStore) SaveCategory(ctx context.Context, cat *Category) error {
	return s.db.WithContext(ctx).Save(cat).Error
}

func (s *gormStore) DeleteCategory(ctx context.Context, cat *Category) error {
	return s.db.WithContext(ctx).Delete(cat).Error
}

func (s *gormStore) GetCategory(ctx context.Context, UUID uuid.UUID) (*Category, error) {
	var cat Category
	if err := s.db.WithContext(ctx).Where("uuid = ?", UUID).First(&cat).Error; err != nil {
		return nil, err
	}
	return &cat, nil
}

func (s *gormStore) GetUserCategories(ctx context.Context, u *users.User) ([]*Category, error) {
	cats := make([]*Category, 0)
	if err := s.db.WithContext(ctx).Where("user_uuid", u.UUID).Find(&cats).Error; err != nil {
		return nil, err
	}
	return cats, nil
}
