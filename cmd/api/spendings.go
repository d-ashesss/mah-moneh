package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type GetSpendingsInput struct {
	Month string `uri:"month" binding:"required,yearmonth"`
}

func (a *App) handleSpendingsGet(c *gin.Context) {
	var input GetSpendingsInput
	if err := c.ShouldBindUri(&input); err != nil {
		c.JSON(http.StatusBadRequest, a.error(err))
		return
	}
	spent, err := a.api.GetUserMonthSpendings(c, a.user(c), input.Month)
	if err != nil {
		panic(fmt.Errorf("failed to get user month spendings: %w", err))
	}
	c.JSON(http.StatusOK, spent)
}
