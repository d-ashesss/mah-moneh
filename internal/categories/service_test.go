package categories_test

import (
	"context"
	"github.com/d-ashesss/mah-moneh/internal/categories"
	mocks "github.com/d-ashesss/mah-moneh/internal/mocks/categories"
	"github.com/d-ashesss/mah-moneh/internal/users"
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
	u := &users.User{}
	cat, err := ts.srv.CreateCategory(ctx, u, "test-cat")
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

func (ts *CategoriesServiceTestSuite) TestGetUserCategories() {
	ctx := context.Background()
	u := &users.User{}
	ts.store.On("GetUserCategories", ctx, u).Return([]*categories.Category{}, nil)
	cats, err := ts.srv.GetUserCategories(ctx, u)
	ts.Require().NoError(err, "Failed to get user categories.")
	ts.Require().NotNil(cats, "Got nil categories.")
}

func TestCategoriesService(t *testing.T) {
	suite.Run(t, new(CategoriesServiceTestSuite))
}
