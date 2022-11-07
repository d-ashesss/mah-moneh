package categories_test

import (
	"context"
	"github.com/d-ashesss/mah-moneh/internal/categories"
	mocks "github.com/d-ashesss/mah-moneh/internal/mocks/categories"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type CategoriesServiceTestSuite struct {
	suite.Suite
	store *mocks.Store
	srv   *categories.Service
}

func (ts *CategoriesServiceTestSuite) SetupTest() {
	ts.store = mocks.NewStore(ts.T())
	ts.srv = categories.NewService(ts.store)
}

func (ts *CategoriesServiceTestSuite) TestCreateCategory() {
	ctx := context.Background()
	ts.store.On("SaveCategory", ctx, mock.AnythingOfType("*categories.Category")).Return(nil)
	cat, err := ts.srv.CreateCategory(ctx, "test-cat", []string{})
	ts.Require().NoError(err, "Failed to create category.")
	ts.Require().NotNil(cat, "Received nil category.")
}

func (ts *CategoriesServiceTestSuite) TestDeleteCategory() {
	ctx := context.Background()
	cat := &categories.Category{}
	ts.store.On("DeleteCategory", ctx, cat).Return(nil)
	err := ts.srv.DeleteCategory(ctx, cat)
	ts.Require().NoError(err, "Failed to delete category.")
}

func TestCategoriesService(t *testing.T) {
	suite.Run(t, new(CategoriesServiceTestSuite))
}
