package categories

import (
	"context"
	"gorm.io/gorm"
)

type Store interface {
	SaveCategory(ctx context.Context, cat *Category) error
	DeleteCategory(ctx context.Context, cat *Category) error
}

type gormStore struct {
	db *gorm.DB
}

func NewGormStore(db *gorm.DB) Store {
	return &gormStore{db: db}
}

func (g *gormStore) SaveCategory(ctx context.Context, cat *Category) error {
	//TODO implement me
	panic("implement SaveCategory")
}

func (g *gormStore) DeleteCategory(ctx context.Context, cat *Category) error {
	//TODO implement me
	panic("implement DeleteCategory")
}
