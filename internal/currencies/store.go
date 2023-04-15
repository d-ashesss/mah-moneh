package currencies

import (
	"errors"
	"github.com/d-ashesss/mah-moneh/internal/accounts"
	"github.com/d-ashesss/mah-moneh/internal/datastore"
	"gorm.io/gorm"
)

// Store is an interface for currencies DB API.
type Store interface {
	// SetRate saves conversion rate into the DB.
	SetRate(base, target accounts.Currency, month string, rate float64) error
	// GetRate retrieves conversion rate from the DB.
	GetRate(base, target accounts.Currency, month string) (*Rate, error)
}

// gormStore is GORM implementation of Store.
type gormStore struct {
	db *gorm.DB
}

func NewGormStore(db *gorm.DB) Store {
	return &gormStore{db: db}
}

func (g *gormStore) SetRate(base, target accounts.Currency, month string, rate float64) error {
	r := &Rate{
		Base:      base,
		Target:    target,
		YearMonth: month,
		Rate:      rate,
	}
	return g.db.Save(r).Error
}

func (g *gormStore) GetRate(base, target accounts.Currency, month string) (*Rate, error) {
	r := &Rate{}
	query := g.db.
		Where("base = ?", base).
		Where("target = ?", target).
		Session(&gorm.Session{})

	err := query.Where("year_month <= ?", month).
		Order("year_month DESC").
		First(r).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = query.Where("year_month > ?", month).
			Order("year_month ASC").
			First(r).Error
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, datastore.ErrRecordNotFound
	}
	if err != nil {
		return nil, err
	}
	return r, nil
}
