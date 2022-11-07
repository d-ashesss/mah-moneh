package categories_test

import (
	"context"
	"fmt"
	"github.com/d-ashesss/mah-moneh/internal/categories"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
	"testing"
)

type CategoriesIntegrationTestSuite struct {
	suite.Suite
	db  *gorm.DB
	srv *categories.Service
}

func (ts *CategoriesIntegrationTestSuite) SetupSuite() {
	dbHost := os.Getenv("POSTGRES_HOST")
	dbPort := os.Getenv("POSTGRES_PORT")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPwd := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	if dbHost == "" {
		ts.T().Skip("No DB configuration provided.")
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s database=%s", dbHost, dbUser, dbPwd, dbName)
	if dbPort != "" {
		dsn = fmt.Sprintf("%s port=%s", dsn, dbPort)
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Discard})
	if err != nil {
		ts.T().Fatalf("Failed to connect to the DB: %s", err)
	}
	ts.db = db
	store := categories.NewGormStore(db)
	ts.srv = categories.NewService(store)
}

func (ts *CategoriesIntegrationTestSuite) SetupTest() {
	_ = ts.db.Migrator().DropTable(&categories.Category{})
	_ = ts.db.AutoMigrate(&categories.Category{})
}

func (ts *CategoriesIntegrationTestSuite) TestSaveCategory() {
	cat, err := ts.srv.CreateCategory(context.Background(), "create-test-category", []string{})
	ts.Require().NoError(err, "Failed to create a category.")
	ts.Require().NotNil(cat, "Received invalid category.")

	foundCat := &categories.Category{}
	err = ts.db.First(foundCat, "uuid = ?", cat.UUID).Error
	ts.Require().NoError(err, "Failed to find created category.")
	ts.NotEmpty(foundCat.UUID)
	ts.Equal(cat.UUID, foundCat.UUID)
}

func (ts *CategoriesIntegrationTestSuite) TestDeleteCategory() {
	cat := categories.NewCategory("delete-test-category", []string{})
	err := ts.db.Save(cat).Error
	ts.Require().NoError(err, "Failed to create testing category.")

	err = ts.srv.DeleteCategory(context.Background(), cat)
	ts.Require().NoError(err, "Failed to delete the category.")

	foundCat := &categories.Category{}
	err = ts.db.First(foundCat, "uuid = ?", cat.UUID).Error
	ts.Require().ErrorIs(err, gorm.ErrRecordNotFound, "Deleted category should not be found.")
}

func TestCategoriesIntegration(t *testing.T) {
	suite.Run(t, new(CategoriesIntegrationTestSuite))
}
