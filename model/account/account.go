package account

import (
	"errors"
	"github.com/d-ashesss/mah-moneh/db"
	"github.com/d-ashesss/mah-moneh/model"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Create saves a new account into the DB.
func Create(u *model.User, name string) (*model.Account, error) {
	acc := model.NewAccount(u, name)
	err := db.DB.Create(acc).Error
	if err != nil {
		return nil, err
	}
	return acc, nil
}

// Get retrieves the account from the DB.
func Get(uuid uuid.UUID) (*model.Account, error) {
	var acc *model.Account
	err := db.DB.First(&acc, "uuid = ?", uuid).Error
	if err != nil {
		return nil, err
	}
	return acc, nil
}

// Update updates account record.
func Update(acc *model.Account) error {
	return db.DB.Where("uuid = ?", acc.UUID).Updates(acc).Error
}

// Delete deletes the account from the DB.
func Delete(acc *model.Account) error {
	return db.DB.Delete(acc).Error
}

// GetAmount retrieves the amount of funds on the account.
func GetAmount(acc *model.Account) (*model.AccountAmount, error) {
	amount := &model.AccountAmount{Account: acc, Amount: 0}
	err := db.DB.First(&amount, "account_uuid = ?", acc.UUID).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return amount, nil
}

// SetAmount sets the amount of funds on the account.
func SetAmount(acc *model.Account, a float64) error {
	amount := &model.AccountAmount{Account: acc, Amount: a}
	return db.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "account_uuid"}},
		DoUpdates: clause.AssignmentColumns([]string{"amount"}),
	}).Create(amount).Error
}
