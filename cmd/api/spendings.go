package main

import (
	"fmt"
	"github.com/d-ashesss/mah-moneh/internal/accounts"
	"github.com/d-ashesss/mah-moneh/internal/categories"
	"github.com/d-ashesss/mah-moneh/internal/spendings"
	"github.com/gin-gonic/gin"
	"net/http"
)

type GetSpendingsInput struct {
	Month string `uri:"month" binding:"required,yearmonth"`
}

type SpendingsResponse map[string]accounts.CurrencyAmounts

func NewSpendingsResponse(spent spendings.Spendings, cats []*categories.Category) SpendingsResponse {
	r := make(SpendingsResponse)
	for _, cat := range cats {
		r[cat.UUID.String()] = spent.GetAmounts(cat)
	}
	r["uncategorized"] = spent.GetAmounts(spendings.Uncategorized)
	r["unaccounted"] = spent.GetAmounts(spendings.Unaccounted)
	r["total"] = spent.GetAmounts(spendings.Total)
	return r
}

func (a *App) handleSpendingsGet(c *gin.Context) {
	var input GetSpendingsInput
	if err := c.ShouldBindUri(&input); err != nil {
		c.JSON(http.StatusBadRequest, a.error(err))
		return
	}
	cats, err := a.api.GetUserCategories(c, a.user(c))
	if err != nil {
		panic(fmt.Errorf("failed to get user categories: %w", err))
	}
	spent, err := a.api.GetUserMonthSpendings(c, a.user(c), input.Month)
	if err != nil {
		panic(fmt.Errorf("failed to get user month spendings: %w", err))
	}
	c.JSON(http.StatusOK, NewSpendingsResponse(spent, cats))
}
