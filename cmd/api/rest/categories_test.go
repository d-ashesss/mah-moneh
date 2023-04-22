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

	ts.Run("create category/invalid name", func() {
		ts.T().Parallel()

		body := bytes.NewBufferString("{}")
		request := NewAuthRequest(user1, "POST", "/categories", body)
		rr := NewRecorder()
		ts.handler.ServeHTTP(rr, request)

		ts.Equal(http.StatusBadRequest, rr.Code)
		response := new(ErrorResponse)
		rr.FromJSON(response)
		ts.Equal("Invalid value of 'Name'", response.Error)
	})

	ts.Run("delete category/invalid id", func() {
		ts.T().Parallel()

		request := NewAuthRequest(user1, "DELETE", "/categories/outsource", nil)
		rr := NewRecorder()
		ts.handler.ServeHTTP(rr, request)

		ts.Equal(http.StatusBadRequest, rr.Code)
		response := new(ErrorResponse)
		rr.FromJSON(response)
		ts.Equal("Invalid value of 'UUID'", response.Error)
	})

	ts.Run("delete category/id not exists", func() {
		ts.T().Parallel()

		randUUID := uuid.Must(uuid.NewV4())
		request := NewAuthRequest(user1, "DELETE", "/categories/"+randUUID.String(), nil)
		rr := NewRecorder()
		ts.handler.ServeHTTP(rr, request)

		ts.Equal(http.StatusNotFound, rr.Code)
		response := new(ErrorResponse)
		rr.FromJSON(response)
		ts.Equal("Not found", response.Error)
	})

	ts.Run("delete category/not owner", func() {
		ts.T().Parallel()

		request := NewAuthRequest(user2, "DELETE", "/categories/"+user1category.UUID.String(), nil)
		rr := NewRecorder()
		ts.handler.ServeHTTP(rr, request)

		ts.Equal(http.StatusNotFound, rr.Code)
		response := new(ErrorResponse)
		rr.FromJSON(response)
		ts.Equal("Not found", response.Error)
	})
}

func (ts *RESTTestSuite) testCreateCategories() {
	tests := []struct {
		Name string
		Body string
		Ref  *uuid.UUID
	}{
		{
			Name: "income",
			Body: `{"name": "income"}`,
			Ref:  &ts.categories.income,
		},
		{
			Name: "groceries",
			Body: `{"name": "groceries"}`,
			Ref:  &ts.categories.income,
		},
		{
			Name: "temp",
			Body: `{"name": "temp"}`,
			Ref:  &ts.categories.income,
		},
	}

	for _, tt := range tests {
		ts.Run(tt.Name, func() {
			ts.T().Parallel()

			body := bytes.NewBufferString(tt.Body)
			request := NewAuthRequest(ts.users.main, "POST", "/categories", body)
			rr := NewRecorder()
			ts.handler.ServeHTTP(rr, request)

			ts.Equal(http.StatusCreated, rr.Code)
			response := new(struct {
				UUID string `json:"uuid"`
			})
			rr.FromJSON(response)
			ts.Require().NotEmptyf(response.UUID, "Received invalid UUID value in response")
			*tt.Ref = uuid.Must(uuid.FromString(response.UUID))
		})
	}
}
