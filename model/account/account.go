package account

import (
	"github.com/d-ashesss/mah-moneh/db"
	"github.com/d-ashesss/mah-moneh/model"
	"github.com/gofrs/uuid"
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
