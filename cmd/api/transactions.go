package main

import (
	"errors"
	"github.com/d-ashesss/mah-moneh/internal/categories"
	"github.com/d-ashesss/mah-moneh/internal/transactions"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"net/http"
)

type CreateTransactionInput struct {
	Month       string  `json:"month" binding:"required,yearmonth"`
	Currency    string  `json:"currency" binding:"required"`
	Amount      float64 `json:"amount" binding:"required"`
	Description string  `json:"description"`
	Category    string  `json:"category_uuid"`
}

type GetTransactionInput struct {
	UUID string `uri:"uuid" binding:"required,uuid"`
}

type GetMonthTransactionsInput struct {
	Month string `uri:"month" binding:"required,yearmonth"`
}

func (a *App) transaction(c *gin.Context) (*transactions.Transaction, error) {
	var input GetTransactionInput
	if err := c.ShouldBindUri(&input); err != nil {
		return nil, err
	}
	tx, err := a.api.GetTransaction(c, uuid.FromStringOrNil(input.UUID))
	if err != nil {
		return nil, err
	}
	if tx.User.UUID != a.user(c).UUID {
		return nil, errors.New(http.StatusText(http.StatusForbidden))
	}
	return tx, nil
}

type TransactionResponse struct {
	UUID        string  `json:"uuid"`
	Month       string  `json:"month"`
	Currency    string  `json:"currency"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
	Category    string  `json:"category_uuid"`
}

func NewTransactionResponse(tx *transactions.Transaction) *TransactionResponse {
	cat := ""
	if tx.Category != nil {
		cat = tx.Category.UUID.String()
	}
	return &TransactionResponse{
		UUID:        tx.UUID.String(),
		Month:       tx.YearMonth,
		Currency:    tx.Currency,
		Amount:      tx.Amount,
		Description: tx.Description,
		Category:    cat,
	}
}

func MapTransactionsResponse(txs []*transactions.Transaction) []*TransactionResponse {
	r := make([]*TransactionResponse, 0, len(txs))
	for _, tx := range txs {
		r = append(r, NewTransactionResponse(tx))
	}
	return r
}

func (a *App) handleTransactionsCreate(c *gin.Context) {
	var input CreateTransactionInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, a.error(err))
		return
	}
	var cat *categories.Category
	if input.Category != "" {
		var err error
		cat, err = a.api.GetCategory(c, uuid.FromStringOrNil(input.Category))
		if err != nil {
			c.JSON(http.StatusBadRequest, a.error(err))
			return
		}
	}
	tx, err := a.api.CreateTransaction(c, a.user(c), input.Month, input.Currency, input.Amount, input.Description, cat)
	if err != nil {
		c.JSON(http.StatusBadRequest, a.error(err))
		return
	}
	c.JSON(http.StatusCreated, NewTransactionResponse(tx))
}

func (a *App) handleTransactionsGet(c *gin.Context) {
	var input GetMonthTransactionsInput
	if err := c.ShouldBindUri(&input); err != nil {
		c.JSON(http.StatusBadRequest, a.error(err))
		return
	}
	txs, err := a.api.GetUserTransactions(c, a.user(c), input.Month)
	if err != nil {
		c.JSON(http.StatusBadRequest, a.error(err))
		return
	}
	c.JSON(http.StatusOK, MapTransactionsResponse(txs))
}

func (a *App) handleTransactionsDelete(c *gin.Context) {
	tx, err := a.transaction(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, a.error(err))
		return
	}
	if err := a.api.DeleteTransaction(c, tx); err != nil {
		c.JSON(http.StatusInternalServerError, a.error(err))
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
