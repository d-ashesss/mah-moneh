//go:build integration

package rest_test

import (
	"bytes"
	"context"
	"github.com/d-ashesss/mah-moneh/internal/users"
	"github.com/gofrs/uuid"
	"net/http"
)

func (ts *RESTTestSuite) testAccountsErrors() {
	user1 := &users.User{UUID: uuid.Must(uuid.NewV4())}
	user2 := &users.User{UUID: uuid.Must(uuid.NewV4())}

	user1account, err := ts.accountsService.CreateAccount(context.Background(), user1, "test account")
	ts.Require().NoErrorf(err, "Failed to create test account")

	tests := []ErrorTest{
		{
			Name:   "create/invalid name",
			Method: "POST",
			Target: "/accounts",
			Auth:   user1,
			Body:   bytes.NewBufferString("{}"),
			Code:   http.StatusBadRequest,
			Error:  "Invalid value of 'Name'",
		},
		{
			Name:   "update/invalid id",
			Method: "PUT",
			Target: "/accounts/wallet",
			Auth:   user1,
			Body:   bytes.NewBufferString("{}"),
			Code:   http.StatusBadRequest,
			Error:  "Invalid value of 'UUID'",
		},
		{
			Name:   "update/id not exists",
			Method: "PUT",
			Target: "/accounts/" + uuid.Must(uuid.NewV4()).String(),
			Auth:   user1,
			Body:   bytes.NewBufferString("{}"),
			Code:   http.StatusNotFound,
			Error:  "Not found",
		},
		{
			Name:   "update/not owner",
			Method: "PUT",
			Target: "/accounts/" + user1account.UUID.String(),
			Auth:   user2,
			Body:   bytes.NewBufferString("{}"),
			Code:   http.StatusNotFound,
			Error:  "Not found",
		},
		{
			Name:   "delete/invalid id",
			Method: "DELETE",
			Target: "/accounts/wallet",
			Auth:   user1,
			Body:   nil,
			Code:   http.StatusBadRequest,
			Error:  "Invalid value of 'UUID'",
		},
		{
			Name:   "delete/id not exists",
			Method: "DELETE",
			Target: "/accounts/" + uuid.Must(uuid.NewV4()).String(),
			Auth:   user1,
			Body:   nil,
			Code:   http.StatusNotFound,
			Error:  "Not found",
		},
		{
			Name:   "delete/not owner",
			Method: "DELETE",
			Target: "/accounts/" + user1account.UUID.String(),
			Auth:   user2,
			Body:   nil,
			Code:   http.StatusNotFound,
			Error:  "Not found",
		},
		{
			Name:   "set amount/invalid month",
			Method: "PUT",
			Target: "/accounts/" + user1account.UUID.String() + "/amounts/201001",
			Auth:   user1,
			Body:   bytes.NewBufferString("{}"),
			Code:   http.StatusBadRequest,
			Error:  "Invalid value of 'Month'",
		},
		{
			Name:   "set amount/invalid currency",
			Method: "PUT",
			Target: "/accounts/" + user1account.UUID.String() + "/amounts/2010-01",
			Auth:   user1,
			Body:   bytes.NewBufferString("{}"),
			Code:   http.StatusBadRequest,
			Error:  "Invalid value of 'Currency'",
		},
		{
			Name:   "set amount/invalid amount data type",
			Method: "PUT",
			Target: "/accounts/" + user1account.UUID.String() + "/amounts/2010-01",
			Auth:   user1,
			Body:   bytes.NewBufferString(`{"currency": "USD", "amount": "100"}`),
			Code:   http.StatusBadRequest,
			Error:  "Invalid request input",
		},
		{
			Name:   "get amount/not owner",
			Method: "GET",
			Target: "/accounts/" + user1account.UUID.String() + "/amounts/2010-01",
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

func (ts *RESTTestSuite) testCreateAccounts() {
	tests := []CreationTest{
		{
			Name: "bank",
			Body: bytes.NewBufferString(`{"name": "bank"}`),
			Ref:  &ts.accounts.bank,
		},
		{
			Name: "cash",
			Body: bytes.NewBufferString(`{"name": "cash"}`),
			Ref:  &ts.accounts.cash,
		},
		{
			Name: "temporary",
			Body: bytes.NewBufferString(`{"name": "tmp"}`),
			Ref:  &ts.accounts.temp,
		},
	}

	for _, tt := range tests {
		ts.testCreate(tt, "/accounts")
	}
}
