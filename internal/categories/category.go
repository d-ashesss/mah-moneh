package categories

import (
	"github.com/d-ashesss/mah-moneh/internal/datastore"
	"github.com/lib/pq"
)

type Category struct {
	datastore.Model
	Name string
	Tags pq.StringArray `gorm:"type:text[]"`
}