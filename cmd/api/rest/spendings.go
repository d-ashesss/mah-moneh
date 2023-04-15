package rest

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

func (i *GetSpendingsInput) Bind(c *gin.Context) error {
	return NewErrBadRequestOrNil(c.ShouldBindUri(i))
}

type SpendingsResponse map[string]accounts.CurrencyAmounts

func NewSpendingsResponse(spent spendings.Spendings, cats []*categories.Category) SpendingsResponse {
	r := make(SpendingsResponse)
	for _, cat := range cats {
		r[cat.UUID.String()] = spent.GetAmounts(cat)
	}
	r["uncategorized"] = spent.GetUncategorized()
	r["unaccounted"] = spent.GetUnaccounted()
	return r
}

func (h *handler) handleSpendingsGet(c *gin.Context) {
	var input GetSpendingsInput
	if err := input.Bind(c); err != nil {
		h.handleError(c, err)
		return
	}
	cats, err := h.categories.GetUserCategories(c, h.user(c))
	if err != nil {
		h.handleError(c, fmt.Errorf("failed to get user categories: %w", err))
		return
	}
	spent, err := h.spendings.GetMonthSpendings(c, h.user(c), input.Month)
	if err != nil {
		h.handleError(c, fmt.Errorf("failed to get user month spendings: %w", err))
		return
	}
	c.JSON(http.StatusOK, NewSpendingsResponse(spent, cats))
}
