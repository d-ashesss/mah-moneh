package spendings_test

import (
	"github.com/d-ashesss/mah-moneh/internal/categories"
	"github.com/d-ashesss/mah-moneh/internal/datastore"
	"github.com/d-ashesss/mah-moneh/internal/spendings"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/suite"
	"testing"
)

type SpendingsTestSuite struct {
	suite.Suite
	category  *categories.Category
	spendings spendings.Spendings
}

func (ts *SpendingsTestSuite) SetupTest() {
	ts.category = &categories.Category{
		Model: datastore.Model{UUID: uuid.NewV5(uuid.Nil, "test")},
	}
	ts.spendings = spendings.NewSpendings([]*categories.Category{ts.category})
}

func (ts *SpendingsTestSuite) TestAddAmount_AddToCategory() {
	ts.spendings.AddAmount(ts.category, "usd", 3)
	ts.spendings.AddAmount(ts.category, "usd", 5.5)
	got := ts.spendings.GetAmount(ts.category, "usd")
	ts.InDelta(8.5, got, 0.001)
	gotTotal := ts.spendings.GetAmount(spendings.Total, "usd")
	ts.InDelta(8.5, gotTotal, 0.001)
}

func (ts *SpendingsTestSuite) TestAddAmount_AddUncategorized() {
	ts.spendings.AddAmount(nil, "usd", 3)
	ts.spendings.AddAmount(nil, "usd", 5.5)
	got := ts.spendings.GetAmount(spendings.Uncategorized, "usd")
	ts.InDelta(8.5, got, 0.001)
	gotTotal := ts.spendings.GetAmount(spendings.Total, "usd")
	ts.InDelta(8.5, gotTotal, 0.001)
}

func (ts *SpendingsTestSuite) TestAddAmount() {
	ts.spendings.AddAmount(ts.category, "usd", 3)
	ts.spendings.AddAmount(nil, "usd", 5.5)
	gotCat := ts.spendings.GetAmount(ts.category, "usd")
	ts.InDelta(3, gotCat, 0.001)
	gotUncat := ts.spendings.GetAmount(spendings.Uncategorized, "usd")
	ts.InDelta(5.5, gotUncat, 0.001)
	gotTotal := ts.spendings.GetAmount(spendings.Total, "usd")
	ts.InDelta(8.5, gotTotal, 0.001)
}

func TestSpendings(t *testing.T) {
	suite.Run(t, new(SpendingsTestSuite))
}
