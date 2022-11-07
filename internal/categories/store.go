package categories

import "context"

type Store interface {
	SaveCategory(ctx context.Context, cat *Category) error
	DeleteCategory(ctx context.Context, cat *Category) error
}
