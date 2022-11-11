package categories

import (
	"github.com/d-ashesss/mah-moneh/internal/datastore"
	"github.com/d-ashesss/mah-moneh/internal/users"
	"github.com/lib/pq"
)

type Category struct {
	datastore.Model
	User *users.User `gorm:"embedded;embeddedPrefix:user_;notNull"`
	Name string
	Tags pq.StringArray `gorm:"type:text[]"`
}

func NewCategory(u *users.User, name string) *Category {
	return &Category{User: u, Name: name}
}
