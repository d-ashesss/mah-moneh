package accounts

import (
	"context"
	"github.com/d-ashesss/mah-moneh/internal/users"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// AccountStore is an interface for accounts DB API.
type AccountStore interface {
	// CreateAccount saves account entity into the DB.
	CreateAccount(ctx context.Context, acc *Account) error
	// UpdateAccount updates existing account entity.
	UpdateAccount(ctx context.Context, acc *Account) error
	// DeleteAccount deletes account entity from the DB.
	DeleteAccount(ctx context.Context, acc *Account) error
	// GetAccount retrieves accounts by its UUID.
	GetAccount(ctx context.Context, UUID uuid.UUID) (*Account, error)
	// GetUserAccounts retrieves all user accounts.
	GetUserAccounts(ctx context.Context, u *users.User) (AccountCollection, error)
	// SetAccountAmount sets the amount of funds on the account.
	SetAccountAmount(ctx context.Context, acc *Account, month string, currency string, amount float64) error
	// GetAccountAmounts retrieves amount of funds for each currency on the account for the specified month.
	GetAccountAmounts(ctx context.Context, acc *Account, month string) (AmountCollection, error)
}

// gormStore is GORM implementation of AccountStore.
type gormStore struct {
	db *gorm.DB
}

// NewGormStore initializes GORM implementation of AccountStore.
func NewGormStore(db *gorm.DB) AccountStore {
	return &gormStore{db: db}
}

func (s *gormStore) CreateAccount(ctx context.Context, acc *Account) error {
	return s.db.WithContext(ctx).Create(acc).Error
}

func (s *gormStore) UpdateAccount(ctx context.Context, acc *Account) error {
	return s.db.WithContext(ctx).Where("uuid = ?", acc.UUID).Updates(acc).Error
}

func (s *gormStore) DeleteAccount(ctx context.Context, acc *Account) error {
	return s.db.WithContext(ctx).Delete(acc).Error
}

func (s *gormStore) GetAccount(ctx context.Context, UUID uuid.UUID) (*Account, error) {
	acc := &Account{}
	if err := s.db.WithContext(ctx).Where("uuid = ?", UUID).First(acc).Error; err != nil {
		return nil, err
	}
	return acc, nil
}

func (s *gormStore) GetUserAccounts(ctx context.Context, u *users.User) (AccountCollection, error) {
	accs := make(AccountCollection, 0)
	if err := s.db.WithContext(ctx).Find(&accs, "user_uuid = ?", u.UUID).Error; err != nil {
		return nil, err
	}
	return accs, nil
}

func (s *gormStore) SetAccountAmount(ctx context.Context, acc *Account, month string, currency string, amount float64) error {
	a := &Amount{Account: acc, YearMonth: month, CurrencyCode: currency, Amount: amount}
	return s.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "account_uuid"},
			{Name: "currency_code"},
			{Name: "year_month"},
		},
		DoUpdates: clause.AssignmentColumns([]string{"amount"}),
	}).Save(a).Error
}

func (s *gormStore) GetAccountAmount(ctx context.Context, acc *Account, month string, currency string) (*Amount, error) {
	amount := &Amount{}
	if err := s.db.WithContext(ctx).
		Where("account_uuid = ?", acc.UUID).
		Where("year_month <= ?", month).
		Where("currency_code = ?", currency).
		Order("year_month desc").
		First(amount).Error; err != nil {
		return nil, err
	}
	return amount, nil
}

func (s *gormStore) GetAccountAmounts(ctx context.Context, acc *Account, month string) (AmountCollection, error) {
	currencies, err := s.GetAccountCurrencies(ctx, acc)
	if err != nil {
		return nil, err
	}
	amounts := make(AmountCollection, 0)
	for _, curr := range currencies {
		amount, err := s.GetAccountAmount(ctx, acc, month, curr)
		if err != nil {
			return nil, err
		}
		amounts = append(amounts, amount)
	}
	return amounts, nil
}

func (s *gormStore) GetAccountCurrencies(ctx context.Context, acc *Account) ([]string, error) {
	var currencies []string
	err := s.db.WithContext(ctx).
		Model(&Amount{}).
		Distinct().
		Where("account_uuid = ?", acc.UUID).
		Pluck("currency_code", &currencies).Error
	if err != nil {
		return nil, err
	}
	return currencies, nil
}
