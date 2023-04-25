//go:build integration

package rest_test

import (
	"github.com/d-ashesss/mah-moneh/internal/users"
	"net/http"
)

func (ts *RESTTestSuite) testGetSpendings() {
	tests := []struct {
		Name     string
		Target   string
		Auth     *users.User
		Expected map[string]map[string]float64
	}{
		{
			Name:   "get main 2009-11 spendings",
			Target: "/spendings/2009-11",
			Auth:   ts.users.main,
			Expected: map[string]map[string]float64{
				ts.categories.income.String():    {},
				ts.categories.groceries.String(): {},
				"uncategorized":                  {},
				"unaccounted":                    {},
			},
		},
		{
			Name:   "get main 2009-12 spendings",
			Target: "/spendings/2009-12",
			Auth:   ts.users.main,
			Expected: map[string]map[string]float64{
				ts.categories.income.String():    {"USD": 2000},
				ts.categories.groceries.String(): {"USD": -300},
				"uncategorized":                  {"USD": -50},
				"unaccounted":                    {"USD": -150},
			},
		},
		{
			Name:   "get main 2010-01 spendings",
			Target: "/spendings/2010-01",
			Auth:   ts.users.main,
			Expected: map[string]map[string]float64{
				ts.categories.income.String(): {
					"USD": 2000,
					"EUR": 500,
				},
				ts.categories.groceries.String(): {"USD": -350},
				"uncategorized":                  {"USD": -200},
				"unaccounted":                    {"USD": -450},
			},
		},
		{
			Name:   "get main 2010-02 spendings",
			Target: "/spendings/2010-02",
			Auth:   ts.users.main,
			Expected: map[string]map[string]float64{
				ts.categories.income.String(): {
					"USD": 1500,
					"EUR": 300,
				},
				ts.categories.groceries.String(): {
					"USD": -250,
					"EUR": -100,
				},
				"uncategorized": {
					"USD": -300,
					"EUR": -200,
				},
				"unaccounted": {"USD": -250},
			},
		},
		{
			Name:   "get main 2010-03 spendings",
			Target: "/spendings/2010-03",
			Auth:   ts.users.main,
			Expected: map[string]map[string]float64{
				ts.categories.income.String(): {
					"USD": 500,
				},
				ts.categories.groceries.String(): {
					"USD": -200,
				},
				"uncategorized": {},
				"unaccounted":   {},
			},
		},
		{
			Name:   "get main 2010-04 spendings",
			Target: "/spendings/2010-04",
			Auth:   ts.users.main,
			Expected: map[string]map[string]float64{
				ts.categories.income.String():    {"USD": 1000},
				ts.categories.groceries.String(): {"USD": -200},
				"uncategorized":                  {},
				"unaccounted":                    {"USD": -800},
			},
		},
		{
			Name:   "get main 2010-05 spendings",
			Target: "/spendings/2010-05",
			Auth:   ts.users.main,
			Expected: map[string]map[string]float64{
				ts.categories.income.String():    {},
				ts.categories.groceries.String(): {},
				"uncategorized":                  {},
				"unaccounted":                    {},
			},
		},

		{
			Name:   "get control 2009-11 spendings",
			Target: "/spendings/2009-11",
			Auth:   ts.users.control,
			Expected: map[string]map[string]float64{
				"uncategorized": {},
				"unaccounted":   {},
			},
		},
	}
	for _, tt := range tests {
		ts.Run(tt.Name, func() {
			request := NewRequest("GET", tt.Target, nil).WithAuth(tt.Auth)
			response := make(map[string]map[string]float64)
			code := ts.ServeJSON(request, &response)

			ts.Equal(http.StatusOK, code)
			ts.Equal(tt.Expected, response)
		})
	}
}
