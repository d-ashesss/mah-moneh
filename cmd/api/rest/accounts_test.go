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

	ts.testSetAccountsAmounts()
}

func (ts *RESTTestSuite) testSetAccountsAmounts() {
	tests := []RequestTest{
		{
			Name:   "cash 2009-12 USD",
			Method: "PUT",
			Target: "/accounts/" + ts.accounts.cash.String() + "/amounts/2009-12",
			Body:   bytes.NewBufferString(`{"currency":"USD","amount":1500}`),
			Auth:   ts.users.main,
			Code:   http.StatusOK,
		},
		{
			Name:   "bank 2010-01 USD",
			Method: "PUT",
			Target: "/accounts/" + ts.accounts.bank.String() + "/amounts/2010-01",
			Body:   bytes.NewBufferString(`{"currency":"USD","amount":2000}`),
			Auth:   ts.users.main,
			Code:   http.StatusOK,
		},
		{
			Name:   "cash 2010-01 USD",
			Method: "PUT",
			Target: "/accounts/" + ts.accounts.cash.String() + "/amounts/2010-01",
			Body:   bytes.NewBufferString(`{"currency":"USD","amount":500}`),
			Auth:   ts.users.main,
			Code:   http.StatusOK,
		},
		{
			Name:   "cash 2010-01 EUR",
			Method: "PUT",
			Target: "/accounts/" + ts.accounts.cash.String() + "/amounts/2010-01",
			Body:   bytes.NewBufferString(`{"currency":"EUR","amount":500}`),
			Auth:   ts.users.main,
			Code:   http.StatusOK,
		},
		{
			Name:   "cash 2010-02 USD",
			Method: "PUT",
			Target: "/accounts/" + ts.accounts.cash.String() + "/amounts/2010-02",
			Body:   bytes.NewBufferString(`{"currency":"USD","amount":1000}`),
			Auth:   ts.users.main,
			Code:   http.StatusOK,
		},
		{
			Name:   "cash 2010-02 EUR",
			Method: "PUT",
			Target: "/accounts/" + ts.accounts.cash.String() + "/amounts/2010-02",
			Body:   bytes.NewBufferString(`{"currency":"EUR","amount":0}`),
			Auth:   ts.users.main,
			Code:   http.StatusOK,
		},
		{
			Name:   "bank 2010-02 USD",
			Method: "PUT",
			Target: "/accounts/" + ts.accounts.bank.String() + "/amounts/2010-02",
			Body:   bytes.NewBufferString(`{"currency":"USD","amount":2200}`),
			Auth:   ts.users.main,
			Code:   http.StatusOK,
		},
		{
			Name:   "bank 2010-02 EUR",
			Method: "PUT",
			Target: "/accounts/" + ts.accounts.bank.String() + "/amounts/2010-02",
			Body:   bytes.NewBufferString(`{"currency":"EUR","amount":500}`),
			Auth:   ts.users.main,
			Code:   http.StatusOK,
		},
		{
			Name:   "temp 2010-02 USD",
			Method: "PUT",
			Target: "/accounts/" + ts.accounts.temp.String() + "/amounts/2010-02",
			Body:   bytes.NewBufferString(`{"currency":"USD","amount":150}`),
			Auth:   ts.users.main,
			Code:   http.StatusOK,
		},
		{
			Name:   "bank 2010-03 USD",
			Method: "PUT",
			Target: "/accounts/" + ts.accounts.bank.String() + "/amounts/2010-03",
			Body:   bytes.NewBufferString(`{"currency":"USD","amount":2500}`),
			Auth:   ts.users.main,
			Code:   http.StatusOK,
		},
	}

	for _, tt := range tests {
		ts.testRequest(tt)
	}
}

func (ts *RESTTestSuite) testDeleteAccounts() {
	tt := RequestTest{
		Name:   "delete account",
		Method: "DELETE",
		Target: "/accounts/" + ts.accounts.temp.String(),
		Body:   nil,
		Auth:   ts.users.main,
		Code:   http.StatusOK,
	}
	ts.testRequest(tt)
}
