package rest

import (
	"fmt"
	"github.com/d-ashesss/mah-moneh/internal/accounts"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"net/http"
)

type CreateAccountInput struct {
	Name string `json:"name" binding:"required"`
}

func (i *CreateAccountInput) Bind(c *gin.Context) error {
	return NewErrBadRequestOrNil(c.ShouldBind(i))
}

type GetAccountInput struct {
	UUID string `uri:"uuid" binding:"required,uuid"`
}

func (i *GetAccountInput) Bind(c *gin.Context) error {
	return NewErrBadRequestOrNil(c.ShouldBindUri(i))
}

func (h *handler) account(c *gin.Context) (*accounts.Account, error) {
	var input GetAccountInput
	if err := input.Bind(c); err != nil {
		return nil, err
	}
	acc, err := h.accounts.GetAccount(c, uuid.FromStringOrNil(input.UUID))
	if err != nil {
		return nil, err
	}
	u := h.user(c)
	if acc.User.UUID != u.UUID {
		return nil, ErrResourceNotFound
	}
	return acc, nil
}

type AccountResponse struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

func NewAccountResponse(acc *accounts.Account) *AccountResponse {
	return &AccountResponse{
		UUID: acc.UUID.String(),
		Name: acc.Name,
	}
}

func NewListAccountsResponse(accs []*accounts.Account) []*AccountResponse {
	r := make([]*AccountResponse, 0, len(accs))
	for _, acc := range accs {
		r = append(r, NewAccountResponse(acc))
	}
	return r
}

func (h *handler) handleAccountsCreate(c *gin.Context) {
	var input CreateAccountInput
	if err := input.Bind(c); err != nil {
		h.handleError(c, err)
		return
	}
	acc, err := h.accounts.CreateAccount(c, h.user(c), input.Name)
	if err != nil {
		h.handleError(c, fmt.Errorf("failed to create account: %w", err))
		return
	}

	c.JSON(http.StatusCreated, NewAccountResponse(acc))
}

func (h *handler) handleAccountsList(c *gin.Context) {
	accs, err := h.accounts.GetUserAccounts(c, h.user(c))
	if err != nil {
		h.handleError(c, fmt.Errorf("failed to get user accounts: %w", err))
		return
	}

	c.JSON(http.StatusOK, NewListAccountsResponse(accs))
}

func (h *handler) handleAccountsUpdate(c *gin.Context) {
	acc, err := h.account(c)
	if err != nil {
		h.handleError(c, fmt.Errorf("failed to find account: %w", err))
		return
	}
	var input CreateAccountInput
	if err := input.Bind(c); err != nil {
		h.handleError(c, err)
		return
	}
	acc.Name = input.Name
	if err := h.accounts.UpdateAccount(c, acc); err != nil {
		h.handleError(c, fmt.Errorf("failed to update account: %w", err))
		return
	}
	c.JSON(http.StatusOK, NewAccountResponse(acc))
}

func (h *handler) handleAccountsDelete(c *gin.Context) {
	acc, err := h.account(c)
	if err != nil {
		h.handleError(c, fmt.Errorf("failed to find account: %w", err))
		return
	}
	if err := h.accounts.DeleteAccount(c, acc); err != nil {
		h.handleError(c, fmt.Errorf("failed to delete account: %w", err))
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

type AccountAmountMonthInput struct {
	Month string `uri:"month" binding:"yearmonth"`
}

func (i *AccountAmountMonthInput) Bind(c *gin.Context) error {
	return NewErrBadRequestOrNil(c.ShouldBindUri(i))
}

type AccountAmountInput struct {
	Currency accounts.Currency `json:"currency" binding:"required"`
	Amount   float64           `json:"amount"`
}

func (i *AccountAmountInput) Bind(c *gin.Context) error {
	return NewErrBadRequestOrNil(c.ShouldBind(i))
}

func (h *handler) handleAccountAmountSet(c *gin.Context) {
	acc, err := h.account(c)
	if err != nil {
		h.handleError(c, fmt.Errorf("failed to find account: %w", err))
		return
	}
	var monthInput AccountAmountMonthInput
	if err := monthInput.Bind(c); err != nil {
		h.handleError(c, err)
		return
	}
	var input AccountAmountInput
	if err := input.Bind(c); err != nil {
		h.handleError(c, err)
		return
	}
	if monthInput.Month == "" {
		if err = h.accounts.SetAccountCurrentAmount(c, acc, input.Currency, input.Amount); err != nil {
			h.handleError(c, fmt.Errorf("failed to set account amount: %w", err))
			return
		}
		c.JSON(http.StatusOK, gin.H{})
		return
	}
	if err = h.accounts.SetAccountAmount(c, acc, monthInput.Month, input.Currency, input.Amount); err != nil {
		h.handleError(c, fmt.Errorf("failed to set account amount: %w", err))
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

func (h *handler) handleAccountAmountGet(c *gin.Context) {
	acc, err := h.account(c)
	if err != nil {
		h.handleError(c, fmt.Errorf("failed to find account: %w", err))
		return
	}
	var monthInput AccountAmountMonthInput
	if err := monthInput.Bind(c); err != nil {
		h.handleError(c, err)
		return
	}
	if monthInput.Month == "" {
		amts, err := h.accounts.GetAccountCurrentAmounts(c, acc)
		if err != nil {
			h.handleError(c, fmt.Errorf("failed to get account amount: %w", err))
			return
		}
		c.JSON(http.StatusOK, amts)
		return
	}
	amts, err := h.accounts.GetAccountAmounts(c, acc, monthInput.Month)
	if err != nil {
		h.handleError(c, fmt.Errorf("failed to get account amount: %w", err))
		return
	}
	c.JSON(http.StatusOK, amts)
}
