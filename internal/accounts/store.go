package accounts

import "gorm.io/gorm"

// AccountStore is an interface for accounts DB API.
type AccountStore interface {
}

// gormStore is GORM implementation of AccountStore.
type gormStore struct {
	db *gorm.DB
}

// NewGormStore initializes GORM implementation of AccountStore.
func NewGormStore(db *gorm.DB) AccountStore {
	return gormStore{db: db}
}
