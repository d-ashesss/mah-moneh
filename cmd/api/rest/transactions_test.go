package rest_test

import (
	"bytes"
	"context"
	"github.com/d-ashesss/mah-moneh/internal/users"
	"github.com/gofrs/uuid"
	"net/http"
)

func (ts *RESTTestSuite) testTransactions() {
	user1 := &users.User{UUID: uuid.Must(uuid.NewV4())}
	user2 := &users.User{UUID: uuid.Must(uuid.NewV4())}

	user1transaction, err := ts.transactionsService.CreateTransaction(context.Background(), user1, "2010-01", "USD", 100, "", nil)
	ts.Require().NoErrorf(err, "Failed to create test transaction")

	ts.Run("create transaction/invalid month", func() {
		ts.T().Parallel()

		body := bytes.NewBufferString(`{"month": "201001"}`)
		request := NewAuthRequest(user1, "POST", "/transactions", body)
		rr := NewRecorder()
		ts.handler.ServeHTTP(rr, request)

		ts.Equal(http.StatusBadRequest, rr.Code)
		response := new(ErrorResponse)
		rr.FromJSON(response)
		ts.Equal("Invalid value of 'Month'", response.Error)
	})

	ts.Run("create transaction/invalid currency", func() {
		ts.T().Parallel()

		body := bytes.NewBufferString(`{"month": "2010-01"}`)
		request := NewAuthRequest(user1, "POST", "/transactions", body)
		rr := NewRecorder()
		ts.handler.ServeHTTP(rr, request)

		ts.Equal(http.StatusBadRequest, rr.Code)
		response := new(ErrorResponse)
		rr.FromJSON(response)
		ts.Equal("Invalid value of 'Currency'", response.Error)
	})

	ts.Run("create transaction/invalid amount", func() {
		ts.T().Parallel()

		body := bytes.NewBufferString(`{"month": "2010-01", "currency": "USD"}`)
		request := NewAuthRequest(user1, "POST", "/transactions", body)
		rr := NewRecorder()
		ts.handler.ServeHTTP(rr, request)

		ts.Equal(http.StatusBadRequest, rr.Code)
		response := new(ErrorResponse)
		rr.FromJSON(response)
		ts.Equal("Invalid value of 'Amount'", response.Error)
	})

	ts.Run("create transaction/invalid amount data type", func() {
		ts.T().Parallel()

		body := bytes.NewBufferString(`{"month": "2010-01", "currency": "USD", "amount": "100"}`)
		request := NewAuthRequest(user1, "POST", "/transactions", body)
		rr := NewRecorder()
		ts.handler.ServeHTTP(rr, request)

		ts.Equal(http.StatusBadRequest, rr.Code)
		response := new(ErrorResponse)
		rr.FromJSON(response)
		ts.Equal("Invalid request input", response.Error)
	})

	ts.Run("create transaction/invalid category", func() {
		ts.T().Parallel()

		body := bytes.NewBufferString(`{"month": "2010-01", "currency": "USD", "amount": 100, "category_uuid": "outsource"}`)
		request := NewAuthRequest(user1, "POST", "/transactions", body)
		rr := NewRecorder()
		ts.handler.ServeHTTP(rr, request)

		ts.Equal(http.StatusNotFound, rr.Code)
		response := new(ErrorResponse)
		rr.FromJSON(response)
		ts.Equal("Not found", response.Error)
	})

	ts.Run("get transactions/invalid month", func() {
		ts.T().Parallel()

		request := NewAuthRequest(user1, "GET", "/transactions/201001", nil)
		rr := NewRecorder()
		ts.handler.ServeHTTP(rr, request)

		ts.Equal(http.StatusBadRequest, rr.Code)
		response := new(ErrorResponse)
		rr.FromJSON(response)
		ts.Equal("Invalid value of 'Month'", response.Error)
	})

	ts.Run("delete transaction/invalid id", func() {
		ts.T().Parallel()

		request := NewAuthRequest(user1, "DELETE", "/transactions/cookies", nil)
		rr := NewRecorder()
		ts.handler.ServeHTTP(rr, request)

		ts.Equal(http.StatusBadRequest, rr.Code)
		response := new(ErrorResponse)
		rr.FromJSON(response)
		ts.Equal("Invalid value of 'UUID'", response.Error)
	})

	ts.Run("delete transaction/id not exists", func() {
		ts.T().Parallel()

		randUUID := uuid.Must(uuid.NewV4())
		request := NewAuthRequest(user1, "DELETE", "/transactions/"+randUUID.String(), nil)
		rr := NewRecorder()
		ts.handler.ServeHTTP(rr, request)

		ts.Equal(http.StatusNotFound, rr.Code)
		response := new(ErrorResponse)
		rr.FromJSON(response)
		ts.Equal("Not found", response.Error)
	})

	ts.Run("delete transaction/not owner", func() {
		ts.T().Parallel()

		request := NewAuthRequest(user2, "DELETE", "/transactions/"+user1transaction.UUID.String(), nil)
		rr := NewRecorder()
		ts.handler.ServeHTTP(rr, request)

		ts.Equal(http.StatusNotFound, rr.Code)
		response := new(ErrorResponse)
		rr.FromJSON(response)
		ts.Equal("Not found", response.Error)
	})
}
