package rest_test

import (
	"bytes"
	"context"
	"fmt"
	"github.com/d-ashesss/mah-moneh/internal/users"
	"github.com/gofrs/uuid"
	"net/http"
)

func (ts *RESTTestSuite) testTransactions() {
	user1 := &users.User{UUID: uuid.Must(uuid.NewV4())}
	user2 := &users.User{UUID: uuid.Must(uuid.NewV4())}

	user1transaction, err := ts.transactionsService.CreateTransaction(context.Background(), user1, "2010-01", "USD", 100, "", nil)
	ts.Require().NoErrorf(err, "Failed to create test transaction")

	tests := []ErrorTest{
		{
			Name:   "create transaction/invalid month",
			Method: "POST",
			Target: "/transactions",
			Auth:   user1,
			Body:   bytes.NewBufferString(`{"month": "201001"}`),
			Code:   http.StatusBadRequest,
			Error:  "Invalid value of 'Month'",
		},
		{
			Name:   "create transaction/invalid currency",
			Method: "POST",
			Target: "/transactions",
			Auth:   user1,
			Body:   bytes.NewBufferString(`{"month": "2010-01"}`),
			Code:   http.StatusBadRequest,
			Error:  "Invalid value of 'Currency'",
		},
		{
			Name:   "create transaction/invalid amount",
			Method: "POST",
			Target: "/transactions",
			Auth:   user1,
			Body:   bytes.NewBufferString(`{"month": "2010-01", "currency": "USD"}`),
			Code:   http.StatusBadRequest,
			Error:  "Invalid value of 'Amount'",
		},
		{
			Name:   "create transaction/invalid amount data type",
			Method: "POST",
			Target: "/transactions",
			Auth:   user1,
			Body:   bytes.NewBufferString(`{"month": "2010-01", "currency": "USD", "amount": "100"}`),
			Code:   http.StatusBadRequest,
			Error:  "Invalid request input",
		},
		{
			Name:   "create transaction/invalid category",
			Method: "POST",
			Target: "/transactions",
			Auth:   user1,
			Body:   bytes.NewBufferString(`{"month": "2010-01", "currency": "USD", "amount": 100, "category_uuid": "outsource"}`),
			Code:   http.StatusNotFound,
			Error:  "Not found",
		},
		{
			Name:   "get transactions/invalid month",
			Method: "GET",
			Target: "/transactions/201001",
			Auth:   user1,
			Body:   nil,
			Code:   http.StatusBadRequest,
			Error:  "Invalid value of 'Month'",
		},
		{
			Name:   "delete transaction/invalid id",
			Method: "DELETE",
			Target: "/transactions/cookies",
			Auth:   user1,
			Body:   nil,
			Code:   http.StatusBadRequest,
			Error:  "Invalid value of 'UUID'",
		},
		{
			Name:   "delete transaction/id not exists",
			Method: "DELETE",
			Target: "/transactions/" + uuid.Must(uuid.NewV4()).String(),
			Auth:   user1,
			Body:   nil,
			Code:   http.StatusNotFound,
			Error:  "Not found",
		},
		{
			Name:   "delete transaction/not owner",
			Method: "DELETE",
			Target: "/transactions/" + user1transaction.UUID.String(),
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

func (ts *RESTTestSuite) testCreateTransactions() {
	tests := []CreationTest{
		// 2009-12
		{
			Name: "2009-12 USD income",
			Body: bytes.NewBufferString(fmt.Sprintf(`{"month": "2009-12","currency": "USD","amount": 2000,"category_uuid": "%s"}`, ts.categories.income)),
			Ref:  nil,
		},
		{
			Name: "2009-12 USD groceries",
			Body: bytes.NewBufferString(fmt.Sprintf(`{"month": "2009-12","currency": "USD","amount": -100,"category_uuid": "%s"}`, ts.categories.groceries)),
			Ref:  nil,
		},
		{
			Name: "2009-12 USD temp",
			Body: bytes.NewBufferString(fmt.Sprintf(`{"month": "2009-12","currency": "USD","amount": -50,"category_uuid": "%s"}`, ts.categories.temp)),
			Ref:  nil,
		},
		{
			Name: "2009-12 USD groceries",
			Body: bytes.NewBufferString(fmt.Sprintf(`{"month": "2009-12","currency": "USD","amount": -200,"category_uuid": "%s"}`, ts.categories.groceries)),
			Ref:  nil,
		},

		// 2010-01
		{
			Name: "2010-01 USD income",
			Body: bytes.NewBufferString(fmt.Sprintf(`{"month": "2010-01","currency": "USD","amount": 2000,"category_uuid": "%s"}`, ts.categories.income)),
			Ref:  nil,
		},
		{
			Name: "2010-01 EUR income",
			Body: bytes.NewBufferString(fmt.Sprintf(`{"month": "2010-01","currency": "EUR","amount": 500,"category_uuid": "%s"}`, ts.categories.income)),
			Ref:  nil,
		},
		{
			Name: "2010-01 USD groceries",
			Body: bytes.NewBufferString(fmt.Sprintf(`{"month": "2010-01","currency": "USD","amount": -200,"category_uuid": "%s"}`, ts.categories.groceries)),
			Ref:  nil,
		},
		{
			Name: "2010-01 USD uncategorized",
			Body: bytes.NewBufferString(fmt.Sprintf(`{"month": "2010-01","currency": "USD","amount": -200}`)),
			Ref:  nil,
		},
		{
			Name: "2010-01 USD groceries",
			Body: bytes.NewBufferString(fmt.Sprintf(`{"month": "2010-01","currency": "USD","amount": -150,"category_uuid": "%s"}`, ts.categories.groceries)),
			Ref:  nil,
		},
		{
			Name: "2010-01 temporary USD groceries",
			Body: bytes.NewBufferString(fmt.Sprintf(`{"month": "2010-01","currency": "USD","amount": -150,"category_uuid": "%s"}`, ts.categories.groceries)),
			Ref:  &ts.transactions.temp,
		},

		// 2010-02
		{
			Name: "2010-02 USD income",
			Body: bytes.NewBufferString(fmt.Sprintf(`{"month": "2010-02","currency": "USD","amount": 1000,"category_uuid": "%s"}`, ts.categories.income)),
			Ref:  nil,
		},
		{
			Name: "2010-02 USD income",
			Body: bytes.NewBufferString(fmt.Sprintf(`{"month": "2010-02","currency": "USD","amount": 500,"category_uuid": "%s"}`, ts.categories.income)),
			Ref:  nil,
		},
		{
			Name: "2010-02 EUR income",
			Body: bytes.NewBufferString(fmt.Sprintf(`{"month": "2010-02","currency": "EUR","amount": 300,"category_uuid": "%s"}`, ts.categories.income)),
			Ref:  nil,
		},
		{
			Name: "2010-02 EUR uncategorized",
			Body: bytes.NewBufferString(fmt.Sprintf(`{"month": "2010-02","currency": "EUR","amount": -200}`)),
			Ref:  nil,
		},
		{
			Name: "2010-02 EUR groceries",
			Body: bytes.NewBufferString(fmt.Sprintf(`{"month": "2010-02","currency": "EUR","amount": -100,"category_uuid": "%s"}`, ts.categories.groceries)),
			Ref:  nil,
		},
		{
			Name: "2010-02 USD uncategorized",
			Body: bytes.NewBufferString(fmt.Sprintf(`{"month": "2010-02","currency": "USD","amount": -300}`)),
			Ref:  nil,
		},
		{
			Name: "2010-02 USD groceries",
			Body: bytes.NewBufferString(fmt.Sprintf(`{"month": "2010-02","currency": "USD","amount": -250,"category_uuid": "%s"}`, ts.categories.groceries)),
			Ref:  nil,
		},

		// 2010-03
		{
			Name: "2010-03 USD income",
			Body: bytes.NewBufferString(fmt.Sprintf(`{"month": "2010-03","currency": "USD","amount": 500,"category_uuid": "%s"}`, ts.categories.income)),
			Ref:  nil,
		},
		{
			Name: "2010-03 USD groceries",
			Body: bytes.NewBufferString(fmt.Sprintf(`{"month": "2010-03","currency": "USD","amount": -200,"category_uuid": "%s"}`, ts.categories.groceries)),
			Ref:  nil,
		},

		// 2010-04
		{
			Name: "2010-04 USD income",
			Body: bytes.NewBufferString(fmt.Sprintf(`{"month": "2010-04","currency": "USD","amount": 1000,"category_uuid": "%s"}`, ts.categories.income)),
			Ref:  nil,
		},
		{
			Name: "2010-04 USD groceries",
			Body: bytes.NewBufferString(fmt.Sprintf(`{"month": "2010-04","currency": "USD","amount": -200,"category_uuid": "%s"}`, ts.categories.groceries)),
			Ref:  nil,
		},
	}

	for _, tt := range tests {
		ts.testCreate(tt, "/transactions")
	}
}
