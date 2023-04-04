//go:build integration

package categories_test

import (
	"context"
	"github.com/d-ashesss/mah-moneh/internal/categories"
	"github.com/d-ashesss/mah-moneh/internal/datastore"
	"github.com/d-ashesss/mah-moneh/internal/users"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"testing"
)

type CategoriesIntegrationTestSuite struct {
	suite.Suite
	db  *gorm.DB
	srv *categories.Service
}

func (ts *CategoriesIntegrationTestSuite) SetupSuite() {
	dbCfg, err := datastore.NewConfig()
	if err != nil {
		ts.T().Fatalf("Invalid database config: %s", err)
	}
	dbCfg.TablePrefix = "cat_test_"
	db, err := datastore.Open(dbCfg)
	if err != nil {
		ts.T().Fatalf("Failed to connect to the DB: %s", err)
	}

	ts.db = db.Session(&gorm.Session{NewDB: true})
	store := categories.NewGormStore(db.Session(&gorm.Session{NewDB: true}))
	ts.srv = categories.NewService(store)

	err = db.Migrator().AutoMigrate(&categories.Category{})
	if err != nil {
		ts.T().Fatalf("Failed to migrate required tables: %s", err)
	}
}

func (ts *CategoriesIntegrationTestSuite) TestSaveCategory() {
	u := ts.createTestingUser()
	cat, err := ts.srv.CreateCategory(context.Background(), u, "create-test-category")
	ts.Require().NoError(err, "Failed to create a category.")
	ts.Require().NotNil(cat, "Received invalid category.")

	foundCat := &categories.Category{}
	err = ts.db.First(foundCat, "uuid = ?", cat.UUID).Error
	ts.Require().NoError(err, "Failed to find created category.")
	ts.NotEmpty(foundCat.UUID)
	ts.Equal(cat.UUID, foundCat.UUID)
}

func (ts *CategoriesIntegrationTestSuite) TestDeleteCategory() {
	u := ts.createTestingUser()
	cat := categories.NewCategory(u, "delete-test-category")
	err := ts.db.Save(cat).Error
	ts.Require().NoError(err, "Failed to create testing category.")

	err = ts.srv.DeleteCategory(context.Background(), cat)
	ts.Require().NoError(err, "Failed to delete the category.")

	foundCat := &categories.Category{}
	err = ts.db.First(foundCat, "uuid = ?", cat.UUID).Error
	ts.Require().ErrorIs(err, gorm.ErrRecordNotFound, "Deleted category should not be found.")
}

func (ts *CategoriesIntegrationTestSuite) TestGetUserCategories() {
	u1 := ts.createTestingUser()
	u2 := ts.createTestingUser()
	var (
		cat  *categories.Category
		cats []*categories.Category
		err  error
	)

	cat = categories.NewCategory(u1, "u1 cat1")
	err = ts.db.Save(cat).Error
	ts.Require().NoError(err, "Failed to create testing category.")

	cat = categories.NewCategory(u1, "u1 cat2")
	err = ts.db.Save(cat).Error
	ts.Require().NoError(err, "Failed to create testing category.")

	cat = categories.NewCategory(u2, "u2 cat1")
	err = ts.db.Save(cat).Error
	ts.Require().NoError(err, "Failed to create testing category.")

	cats, err = ts.srv.GetUserCategories(context.Background(), u1)
	ts.Require().NoError(err, "Failed to get user's categories.")
	ts.Len(cats, 2, "Got invalid number of categories.")

	cats, err = ts.srv.GetUserCategories(context.Background(), u2)
	ts.Require().NoError(err, "Failed to get user's categories.")
	ts.Len(cats, 1, "Got invalid number of categories.")
}

func (ts *CategoriesIntegrationTestSuite) createTestingUser() *users.User {
	ts.T().Helper()
	UUID, _ := uuid.NewV4()
	return &users.User{UUID: UUID}
}

func TestCategoriesIntegration(t *testing.T) {
	suite.Run(t, new(CategoriesIntegrationTestSuite))
}
