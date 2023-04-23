package rest_test

import (
	"bytes"
	"context"
	"github.com/d-ashesss/mah-moneh/internal/users"
	"github.com/gofrs/uuid"
	"net/http"
)

func (ts *RESTTestSuite) testCategoriesErrors() {
	user1 := &users.User{UUID: uuid.Must(uuid.NewV4())}
	user2 := &users.User{UUID: uuid.Must(uuid.NewV4())}

	user1category, err := ts.categoriesService.CreateCategory(context.Background(), user1, "test category")
	ts.Require().NoErrorf(err, "Failed to create test category")

	tests := []ErrorTest{
		{
			Name:   "create category/invalid name",
			Method: "POST",
			Target: "/categories",
			Auth:   user1,
			Body:   bytes.NewBufferString(`{}`),
			Code:   http.StatusBadRequest,
			Error:  "Invalid value of 'Name'",
		},
		{
			Name:   "delete category/invalid id",
			Method: "DELETE",
			Target: "/categories/outsource",
			Auth:   user1,
			Body:   nil,
			Code:   http.StatusBadRequest,
			Error:  "Invalid value of 'UUID'",
		},
		{
			Name:   "delete category/id not exists",
			Method: "DELETE",
			Target: "/categories/" + uuid.Must(uuid.NewV4()).String(),
			Auth:   user1,
			Body:   nil,
			Code:   http.StatusNotFound,
			Error:  "Not found",
		},
		{
			Name:   "delete category/not owner",
			Method: "DELETE",
			Target: "/categories/" + user1category.UUID.String(),
			Auth:   user2,
			Body:   nil,
			Code:   http.StatusNotFound,
			Error:  "Not found",
		},
	}

	for _, tt := range tests {
		ts.testError(tt)
	}
}

func (ts *RESTTestSuite) testCreateCategories() {
	tests := []CreationTest{
		{
			Name: "income",
			Body: bytes.NewBufferString(`{"name": "income"}`),
			Ref:  &ts.categories.income,
		},
		{
			Name: "groceries",
			Body: bytes.NewBufferString(`{"name": "groceries"}`),
			Ref:  &ts.categories.groceries,
		},
		{
			Name: "temp",
			Body: bytes.NewBufferString(`{"name": "temp"}`),
			Ref:  &ts.categories.temp,
		},
	}

	for _, tt := range tests {
		ts.testCreate(tt, "/categories")
	}
}
