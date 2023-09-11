package rest

import (
	"fmt"
	"github.com/d-ashesss/mah-moneh/internal/accounts"
	"github.com/d-ashesss/mah-moneh/internal/categories"
	"github.com/d-ashesss/mah-moneh/internal/transactions"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"net/http"
)

type CreateTransactionInput struct {
	Month        string            `json:"month" binding:"required,yearmonth"`
	Currency     accounts.Currency `json:"currency" binding:"required"`
	Amount       float64           `json:"amount" binding:"required"`
	Description  string            `json:"description"`
	CategoryUUID *string           `json:"category_uuid"`
}

func (i *CreateTransactionInput) Bind(c *gin.Context) error {
	return NewErrBadRequestOrNil(c.ShouldBind(i))
}

type GetTransactionInput struct {
	UUID string `uri:"uuid" binding:"required,uuid"`
}

func (i *GetTransactionInput) Bind(c *gin.Context) error {
	return NewErrBadRequestOrNil(c.ShouldBindUri(i))
}

type GetMonthTransactionsInput struct {
	Month string `uri:"month" binding:"required,yearmonth"`
}

func (i *GetMonthTransactionsInput) Bind(c *gin.Context) error {
	return NewErrBadRequestOrNil(c.ShouldBindUri(i))
}

func (h *handler) transactionCategory(c *gin.Context, i CreateTransactionInput) (*categories.Category, error) {
	if i.CategoryUUID != nil {
		return h.categories.GetCategory(c, uuid.FromStringOrNil(*i.CategoryUUID))
	}
	return nil, nil
}

func (h *handler) transaction(c *gin.Context) (*transactions.Transaction, error) {
	var input GetTransactionInput
	if err := input.Bind(c); err != nil {
		return nil, err
	}
	tx, err := h.transactions.GetTransaction(c, uuid.FromStringOrNil(input.UUID))
	if err != nil {
		return nil, err
	}
	if tx.User.UUID != h.user(c).UUID {
		return nil, ErrResourceNotFound
	}
	return tx, nil
}

type TransactionResponse struct {
	UUID        string            `json:"uuid"`
	Month       string            `json:"month"`
	Currency    accounts.Currency `json:"currency"`
	Amount      float64           `json:"amount"`
	Description string            `json:"description"`
	Category    string            `json:"category_uuid"`
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

func NewListTransactionsResponse(txs []*transactions.Transaction) []*TransactionResponse {
	r := make([]*TransactionResponse, 0, len(txs))
	for _, tx := range txs {
		r = append(r, NewTransactionResponse(tx))
	}
	return r
}

func (h *handler) handleTransactionsCreate(c *gin.Context) {
	var input CreateTransactionInput
	if err := input.Bind(c); err != nil {
		h.handleError(c, err)
		return
	}
	cat, err := h.transactionCategory(c, input)
	if err != nil {
		h.handleError(c, err)
		return
	}
	tx, err := h.transactions.CreateTransaction(c, h.user(c), input.Month, input.Currency, input.Amount, input.Description, cat)
	if err != nil {
		h.handleError(c, fmt.Errorf("failed to create transaction: %w", err))
		return
	}
	c.JSON(http.StatusCreated, NewTransactionResponse(tx))
}

func (h *handler) handleTransactionsList(c *gin.Context) {
	var input GetMonthTransactionsInput
	if err := input.Bind(c); err != nil {
		h.handleError(c, err)
		return
	}
	txs, err := h.transactions.GetUserTransactions(c, h.user(c), input.Month)
	if err != nil {
		h.handleError(c, fmt.Errorf("failed to get user transactions: %w", err))
		return
	}
	c.JSON(http.StatusOK, NewListTransactionsResponse(txs))
}

func (h *handler) handleTransactionsDelete(c *gin.Context) {
	tx, err := h.transaction(c)
	if err != nil {
		h.handleError(c, fmt.Errorf("failed to find transaction: %w", err))
		return
	}
	if err := h.transactions.DeleteTransaction(c, tx); err != nil {
		h.handleError(c, fmt.Errorf("failed to delete transaction: %w", err))
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
