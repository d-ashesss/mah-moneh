//go:build integration

package rest_test

import (
	"bytes"
	"context"
	"github.com/gofrs/uuid"
	"net/http"
)

func (ts *RESTTestSuite) testCategoriesErrors() {
	auth1 := ts.NewAuth()
	auth2 := ts.NewAuth()

	user1category, err := ts.categoriesService.CreateCategory(context.Background(), auth1.user, "test category")
	ts.Require().NoErrorf(err, "Failed to create test category")

	tests := []ErrorTest{
		{
			Name:   "create category/invalid name",
			Method: "POST",
			Target: "/categories",
			Auth:   auth1,
			Body:   bytes.NewBufferString(`{}`),
			Code:   http.StatusBadRequest,
			Error:  "Invalid value of 'Name'",
		},
		{
			Name:   "delete category/invalid id",
			Method: "DELETE",
			Target: "/categories/outsource",
			Auth:   auth1,
			Body:   nil,
			Code:   http.StatusBadRequest,
			Error:  "Invalid value of 'UUID'",
		},
		{
			Name:   "delete category/id not exists",
			Method: "DELETE",
			Target: "/categories/" + uuid.Must(uuid.NewV4()).String(),
			Auth:   auth1,
			Body:   nil,
			Code:   http.StatusNotFound,
			Error:  "Not found",
		},
		{
			Name:   "delete category/not owner",
			Method: "DELETE",
			Target: "/categories/" + user1category.UUID.String(),
			Auth:   auth2,
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

func (ts *RESTTestSuite) testDeleteCategories() {
	tt := RequestTest{
		Name:   "delete category",
		Method: "DELETE",
		Target: "/categories/" + ts.categories.temp.String(),
		Body:   nil,
		Auth:   ts.users.main,
		Code:   http.StatusOK,
	}
	ts.testRequest(tt)
}

func (ts *RESTTestSuite) testGetCategories() {
	tests := []CountTest{
		{
			Name:   "get main categories",
			Target: "/categories",
			Auth:   ts.users.main,
			Count:  2,
		},
		{
			Name:   "get control categories",
			Target: "/categories",
			Auth:   ts.users.control,
			Count:  0,
		},
	}
	for _, tt := range tests {
		ts.testCount(tt)
	}
}
