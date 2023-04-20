//go:build integration

package rest_test

import (
	"bytes"
	"context"
	"github.com/d-ashesss/mah-moneh/internal/users"
	"github.com/gofrs/uuid"
	"net/http"
)

func (ts *RESTTestSuite) testAccounts() {
	user1 := &users.User{UUID: uuid.Must(uuid.NewV4())}
	user2 := &users.User{UUID: uuid.Must(uuid.NewV4())}

	user1account, err := ts.accountsService.CreateAccount(context.Background(), user1, "test account")
	ts.Require().NoErrorf(err, "Failed to create test account")

	ts.Run("create/invalid name", func() {
		ts.T().Parallel()

		body := bytes.NewBufferString(`{}`)
		request := NewAuthRequest(user1, "POST", "/accounts", body)
		rr := NewRecorder()
		ts.handler.ServeHTTP(rr, request)

		ts.Equal(http.StatusBadRequest, rr.Code)
		response := new(ErrorResponse)
		rr.FromJSON(&response)
		ts.Equal("Invalid value of 'Name'", response.Error)
	})

	ts.Run("update/invalid id", func() {
		ts.T().Parallel()

		body := bytes.NewBufferString(`{}`)
		request := NewAuthRequest(user1, "PUT", "/accounts/wallet", body)
		rr := NewRecorder()
		ts.handler.ServeHTTP(rr, request)

		ts.Equal(http.StatusBadRequest, rr.Code)
		response := new(ErrorResponse)
		rr.FromJSON(&response)
		ts.Equal("Invalid value of 'UUID'", response.Error)
	})

	ts.Run("update/id not exists", func() {
		ts.T().Parallel()

		body := bytes.NewBufferString(`{}`)
		randUUID := uuid.Must(uuid.NewV4())
		request := NewAuthRequest(user1, "PUT", "/accounts/"+randUUID.String(), body)

		rr := NewRecorder()
		ts.handler.ServeHTTP(rr, request)

		ts.Equal(http.StatusNotFound, rr.Code)
		response := new(ErrorResponse)
		rr.FromJSON(&response)
		ts.Equal("Not found", response.Error)
	})

	ts.Run("update/not owner", func() {
		ts.T().Parallel()

		body := bytes.NewBufferString(`{}`)
		request := NewAuthRequest(user2, "PUT", "/accounts/"+user1account.UUID.String(), body)

		rr := NewRecorder()
		ts.handler.ServeHTTP(rr, request)

		ts.Equal(http.StatusNotFound, rr.Code)
		response := new(ErrorResponse)
		rr.FromJSON(&response)
		ts.Equal("Not found", response.Error)
	})

	ts.Run("delete/invalid id", func() {
		ts.T().Parallel()

		request := NewAuthRequest(user1, "DELETE", "/accounts/wallet", nil)
		rr := NewRecorder()
		ts.handler.ServeHTTP(rr, request)

		ts.Equal(http.StatusBadRequest, rr.Code)
		response := new(ErrorResponse)
		rr.FromJSON(&response)
		ts.Equal("Invalid value of 'UUID'", response.Error)
	})

	ts.Run("delete/id not exists", func() {
		ts.T().Parallel()

		randUUID := uuid.Must(uuid.NewV4())
		request := NewAuthRequest(user1, "DELETE", "/accounts/"+randUUID.String(), nil)

		rr := NewRecorder()
		ts.handler.ServeHTTP(rr, request)

		ts.Equal(http.StatusNotFound, rr.Code)
		response := new(ErrorResponse)
		rr.FromJSON(&response)
		ts.Equal("Not found", response.Error)
	})

	ts.Run("delete/not owner", func() {
		ts.T().Parallel()

		request := NewAuthRequest(user2, "DELETE", "/accounts/"+user1account.UUID.String(), nil)

		rr := NewRecorder()
		ts.handler.ServeHTTP(rr, request)

		ts.Equal(http.StatusNotFound, rr.Code)
		response := new(ErrorResponse)
		rr.FromJSON(&response)
		ts.Equal("Not found", response.Error)
	})

	ts.Run("set amount/invalid month", func() {
		ts.T().Parallel()

		body := bytes.NewBufferString(`{}`)
		request := NewAuthRequest(user1, "PUT", "/accounts/"+user1account.UUID.String()+"/amounts/201001", body)
		rr := NewRecorder()
		ts.handler.ServeHTTP(rr, request)

		ts.Equal(http.StatusBadRequest, rr.Code)
		response := new(ErrorResponse)
		rr.FromJSON(&response)
		ts.Equal("Invalid value of 'Month'", response.Error)
	})

	ts.Run("set amount/invalid currency", func() {
		ts.T().Parallel()

		body := bytes.NewBufferString(`{}`)
		request := NewAuthRequest(user1, "PUT", "/accounts/"+user1account.UUID.String()+"/amounts/2010-01", body)
		rr := NewRecorder()
		ts.handler.ServeHTTP(rr, request)

		ts.Equal(http.StatusBadRequest, rr.Code)
		response := new(ErrorResponse)
		rr.FromJSON(&response)
		ts.Equal("Invalid value of 'Currency'", response.Error)
	})

	ts.Run("set amount/invalid amount data type", func() {
		ts.T().Parallel()

		body := bytes.NewBufferString(`{"currency": "USD", "amount": "100"}`)
		request := NewAuthRequest(user1, "PUT", "/accounts/"+user1account.UUID.String()+"/amounts/2010-01", body)
		rr := NewRecorder()
		ts.handler.ServeHTTP(rr, request)

		ts.Equal(http.StatusBadRequest, rr.Code)
		response := new(ErrorResponse)
		rr.FromJSON(&response)
		ts.Equal("Invalid request input", response.Error)
	})

	ts.Run("get amount/not owner", func() {
		ts.T().Parallel()

		request := NewAuthRequest(user2, "GET", "/accounts/"+user1account.UUID.String()+"/amounts/2010-01", nil)
		rr := NewRecorder()
		ts.handler.ServeHTTP(rr, request)

		ts.Equal(http.StatusNotFound, rr.Code)
		response := new(ErrorResponse)
		rr.FromJSON(&response)
		ts.Equal("Not found", response.Error)
	})
}
