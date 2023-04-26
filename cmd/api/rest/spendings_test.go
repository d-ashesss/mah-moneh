//go:build integration

package rest_test

import "fmt"

func (ts *RESTTestSuite) testGetSpendings() {
	tests := []JSONTest{
		{
			Name:   "get main 2009-11 spendings",
			Target: "/spendings/2009-11",
			Auth:   ts.users.main,
			Expected: fmt.Sprintf(`{
				"%s":            {},
				"%s":            {},
				"uncategorized": {},
				"unaccounted":   {}
			}`, ts.categories.income, ts.categories.groceries),
		},
		{
			Name:   "get main 2009-12 spendings",
			Target: "/spendings/2009-12",
			Auth:   ts.users.main,
			Expected: fmt.Sprintf(`{
				"%s":            {"USD": 2000},
				"%s":            {"USD": -300},
				"uncategorized": {"USD": -50},
				"unaccounted":   {"USD": -150}
			}`, ts.categories.income, ts.categories.groceries),
		},
		{
			Name:   "get main 2010-01 spendings",
			Target: "/spendings/2010-01",
			Auth:   ts.users.main,
			Expected: fmt.Sprintf(`{
				"%s":            {"USD": 2000, "EUR": 500},
				"%s":            {"USD": -350},
				"uncategorized": {"USD": -200},
				"unaccounted":   {"USD": -450}
			}`, ts.categories.income, ts.categories.groceries),
		},
		{
			Name:   "get main 2010-02 spendings",
			Target: "/spendings/2010-02",
			Auth:   ts.users.main,
			Expected: fmt.Sprintf(`{
				"%s":            {"USD": 1500, "EUR": 300},
				"%s":            {"USD": -250, "EUR": -100},
				"uncategorized": {"USD": -300, "EUR": -200},
				"unaccounted":   {"USD": -250}
			}`, ts.categories.income, ts.categories.groceries),
		},
		{
			Name:   "get main 2010-03 spendings",
			Target: "/spendings/2010-03",
			Auth:   ts.users.main,
			Expected: fmt.Sprintf(`{
				"%s":            {"USD": 500},
				"%s":            {"USD": -200},
				"uncategorized": {},
				"unaccounted":   {}
			}`, ts.categories.income, ts.categories.groceries),
		},
		{
			Name:   "get main 2010-04 spendings",
			Target: "/spendings/2010-04",
			Auth:   ts.users.main,
			Expected: fmt.Sprintf(`{
				"%s":            {"USD": 1000},
				"%s":            {"USD": -200},
				"uncategorized": {},
				"unaccounted":   {"USD": -800}
			}`, ts.categories.income, ts.categories.groceries),
		},
		{
			Name:   "get main 2010-05 spendings",
			Target: "/spendings/2010-05",
			Auth:   ts.users.main,
			Expected: fmt.Sprintf(`{
				"%s":            {},
				"%s":            {},
				"uncategorized": {},
				"unaccounted":   {}
			}`, ts.categories.income, ts.categories.groceries),
		},

		{
			Name:   "get control 2009-11 spendings",
			Target: "/spendings/2009-11",
			Auth:   ts.users.control,
			Expected: `{
				"uncategorized": {},
				"unaccounted":   {}
			}`,
		},
	}
	for _, tt := range tests {
		ts.testJSON(tt)
	}
}
