package model_test

import (
	"github.com/d-ashesss/mah-moneh/db"
	"github.com/d-ashesss/mah-moneh/model"
	"github.com/d-ashesss/mah-moneh/model/account"
	"github.com/d-ashesss/mah-moneh/model/capital"
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
			u := &model.User{UUID: uuid.FromStringOrNil("abbad7a4-1922-4f11-aefd-f55b563816ce")}
			newacc, err := account.Create(u, "test")
			if !assert.NoError(t, err, "Failed to create account") {
				return
			}
			acc, err := account.Get(newacc.UUID)
			if assert.NoError(t, err) {
				assert.Equal(t, newacc.UUID, acc.UUID)
			}
		})

		t.Run("GetAll", func(t *testing.T) {
			u := &model.User{UUID: uuid.FromStringOrNil("9de1799e-e5be-40a2-9715-e7ba513426e4")}
			accs, err := account.GetAll(u)
			if !assert.NoError(t, err, "Failed to get accounts") {
				return
			}
			assert.NotNil(t, accs)
			assert.Empty(t, accs)

			_, err = account.Create(u, "test")
			if !assert.NoError(t, err) {
				return
			}

			accs, err = account.GetAll(u)
			if !assert.NoError(t, err, "Failed to get accounts") {
				return
			}
			assert.Len(t, accs, 1)
		})

		t.Run("Update", func(t *testing.T) {
			u := &model.User{UUID: uuid.FromStringOrNil("a6ab6796-08fb-4fb8-bd15-faa63e5378ee")}
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
			u := &model.User{UUID: uuid.FromStringOrNil("d5870983-d2ec-4f21-ad3a-37ff41163912")}
			acc, err := account.Create(u, "test")
			if !assert.NoError(t, err, "Failed to create account") {
				return
			}
			err = account.Delete(acc)
			if !assert.NoError(t, err, "Failed to delete account") {
				return
			}
			_, err = account.Get(acc.UUID)
			assert.ErrorContains(t, err, "not found", "Deleted account should not be found")
		})

		t.Run("GetAmount", func(t *testing.T) {
			u := &model.User{UUID: uuid.FromStringOrNil("88d962b4-63eb-472b-8f61-63af9c0f3e19")}
			acc, err := account.Create(u, "test")
			if !assert.NoError(t, err, "Failed to create account") {
				return
			}
			a, err := account.GetAmount(acc)
			if !assert.NoError(t, err, "Failed to get amount") {
				return
			}
			assert.Equal(t, 0., a.Amount, "Unexpected amount value")
		})

		t.Run("SetAmount", func(t *testing.T) {
			u := &model.User{UUID: uuid.FromStringOrNil("58f3fb41-651d-4b14-b409-ef9f6c4f8fe9")}
			acc, err := account.Create(u, "test")
			if !assert.NoError(t, err, "Failed to create account") {
				return
			}
			err = account.SetAmount(acc, 10)
			if !assert.NoError(t, err, "Failed to set amount") {
				return
			}
			a, err := account.GetAmount(acc)
			if !assert.NoError(t, err, "Failed to get amount") {
				return
			}
			if !assert.Equal(t, float64(10), a.Amount, "Unexpected amount value") {
				return
			}
			err = account.SetAmount(acc, 20)
			if !assert.NoError(t, err, "Failed to set amount") {
				return
			}
			a, err = account.GetAmount(acc)
			if !assert.NoError(t, err, "Failed to get amount") {
				return
			}
			if !assert.Equal(t, float64(20), a.Amount, "Unexpected amount value") {
				return
			}
		})
	})

	t.Run("Capital", func(t *testing.T) {
		assertCapital := func(t *testing.T, u *model.User, amount float64) {
			t.Helper()
			c, err := capital.Get(u)
			if !assert.NoError(t, err, "Failed to get capital") {
				return
			}
			if !assert.Equal(t, amount, c.Amount, "Unexpected capital amount value") {
				return
			}
		}

		t.Run("Get", func(t *testing.T) {
			u := &model.User{UUID: uuid.FromStringOrNil("c748a4ea-124e-4832-bfeb-dd529f4fbb9c")}
			acc1, err := account.Create(u, "first")
			if !assert.NoError(t, err, "Failed to create account") {
				return
			}
			err = account.SetAmount(acc1, 10)
			if !assert.NoError(t, err, "Failed to set amount") {
				return
			}
			assertCapital(t, u, 10)

			acc2, err := account.Create(u, "second")
			if !assert.NoError(t, err, "Failed to create account") {
				return
			}
			err = account.SetAmount(acc2, 15)
			if !assert.NoError(t, err, "Failed to set amount") {
				return
			}
			assertCapital(t, u, 25)
		})
	})
}

func createDb(t *testing.T) {
	t.Helper()
	var err error

	err = db.DB.Migrator().DropTable(&model.AccountAmount{})
	if err != nil {
		t.Fatalf("[main] migrate: %v", err)
	}

	err = db.DB.Migrator().DropTable(&model.Account{})
	if err != nil {
		t.Fatalf("[main] migrate: %v", err)
	}

	err = db.DB.Migrator().CreateTable(&model.Account{})
	if err != nil {
		t.Fatalf("[main] migrate: %v", err)
	}

	err = db.DB.Migrator().CreateTable(&model.AccountAmount{})
	if err != nil {
		t.Fatalf("[main] migrate: %v", err)
	}
}
