package rest_test

import (
	"bytes"
	"context"
	"github.com/d-ashesss/mah-moneh/internal/users"
	"github.com/gofrs/uuid"
	"net/http"
)

func (ts *RESTTestSuite) testCategories() {
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
