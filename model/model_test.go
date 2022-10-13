package model_test

import (
	"github.com/d-ashesss/mah-moneh/db"
	"github.com/d-ashesss/mah-moneh/model"
	"github.com/d-ashesss/mah-moneh/model/account"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestModel(t *testing.T) {
	if err := db.Connect("172.17.0.1", "dash", "loc", "test"); err != nil {
		t.Fatalf("connect: %v", err)
	}
	createDb(t)

	t.Run("Account", func(t *testing.T) {
		t.Run("Create", func(t *testing.T) {
			u := &model.User{UUID: uuid.FromStringOrNil("3d8d0573-1240-46c3-bacb-3c172047eb6a")}
			acc, err := account.Create(u, "test")
			if assert.NoError(t, err) {
				assert.NotEqual(t, uuid.Nil, acc.UUID, "Expected UUID to be set")
			}
		})

		t.Run("Get", func(t *testing.T) {
			u := &model.User{UUID: uuid.FromStringOrNil("3d8d0573-1240-46c3-bacb-3c172047eb6a")}
			newacc, err := account.Create(u, "test")
			if !assert.NoError(t, err, "Failed to create account") {
				return
			}
			acc, err := account.Get(newacc.UUID)
			if assert.NoError(t, err) {
				assert.Equal(t, newacc.UUID, acc.UUID)
			}
		})

		t.Run("Update", func(t *testing.T) {
			u := &model.User{UUID: uuid.FromStringOrNil("3d8d0573-1240-46c3-bacb-3c172047eb6a")}
			newacc, err := account.Create(u, "test")
			if !assert.NoError(t, err, "Failed to create account") {
				return
			}
			newacc.Name = "new test"
			err = account.Update(newacc)
			if !assert.NoError(t, err) {
				return
			}
			acc, err := account.Get(newacc.UUID)
			if !assert.NoError(t, err) {
				return
			}
			assert.Equal(t, newacc.Name, acc.Name, "Name was not changed")
		})

		t.Run("Delete", func(t *testing.T) {
			u := &model.User{UUID: uuid.FromStringOrNil("3d8d0573-1240-46c3-bacb-3c172047eb6a")}
			newacc, err := account.Create(u, "test")
			if !assert.NoError(t, err, "Failed to create account") {
				return
			}
			err = account.Delete(newacc)
			if !assert.NoError(t, err, "Failed to delete account") {
				return
			}
			_, err = account.Get(newacc.UUID)
			assert.ErrorContains(t, err, "not found", "Deleted account should not be found")
		})
	})
}

func createDb(t *testing.T) {
	t.Helper()
	err := db.DB.Migrator().DropTable(&model.Account{})
	if err != nil {
		t.Fatalf("[main] migrate: %v", err)
	}
	err = db.DB.Migrator().CreateTable(&model.Account{})
	if err != nil {
		t.Fatalf("[main] migrate: %v", err)
	}
}
