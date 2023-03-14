package main

import (
	"errors"
	"fmt"
	"github.com/d-ashesss/mah-moneh/internal/accounts"
	"github.com/d-ashesss/mah-moneh/internal/api"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"net/http"
	"time"
)

type CreateAccountInput struct {
	Name string `json:"name" binding:"required"`
}

type GetAccountInput struct {
	UUID string `uri:"uuid" binding:"required,uuid"`
}

func (a *App) account(c *gin.Context) (*accounts.Account, error) {
	var input GetAccountInput
	if err := c.ShouldBindUri(&input); err != nil {
		return nil, err
	}
	acc, err := a.api.GetAccount(c, uuid.FromStringOrNil(input.UUID))
	if err != nil {
		return nil, err
	}
	u := a.user(c)
	if acc.User.UUID != u.UUID {
		return nil, api.ErrResourceNotFound
	}
	return acc, nil
}

type AccountResponse struct {
	UUID      string `json:"uuid"`
	CreatedAt string `json:"created_at"`
	Name      string `json:"name"`
}

func NewAccountResponse(acc *accounts.Account) *AccountResponse {
	return &AccountResponse{
		UUID:      acc.UUID.String(),
		CreatedAt: acc.CreatedAt.Format(time.DateTime),
		Name:      acc.Name,
	}
}

func MapAccountsResponse(accs []*accounts.Account) []*AccountResponse {
	r := make([]*AccountResponse, 0, len(accs))
	for _, acc := range accs {
		r = append(r, NewAccountResponse(acc))
	}
	return r
}

func (a *App) handleAccountsCreate(c *gin.Context) {
	var input CreateAccountInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, a.error(err))
		return
	}
	acc, err := a.api.CreateAccount(c, a.user(c), input.Name)
	if err != nil {
		panic(fmt.Errorf("failed to create account: %w", err))
	}

	c.JSON(http.StatusCreated, NewAccountResponse(acc))
}

func (a *App) handleAccountsGet(c *gin.Context) {
	accs, err := a.api.GetUserAccounts(c, a.user(c))
	if err != nil {
		panic(fmt.Errorf("failed to get user accounts: %w", err))
	}

	c.JSON(http.StatusOK, MapAccountsResponse(accs))
}

func (a *App) handleAccountsUpdate(c *gin.Context) {
	acc, err := a.account(c)
	if errors.Is(err, api.ErrResourceNotFound) {
		c.JSON(http.StatusNotFound, a.error("Account not found"))
		return
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, a.error(err))
		return
	}
	var input CreateAccountInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, a.error(err))
		return
	}
	acc.Name = input.Name
	if err := a.api.UpdateAccount(c, acc); err != nil {
		panic(fmt.Errorf("failed to update account: %w", err))
	}
	c.JSON(http.StatusOK, NewAccountResponse(acc))
}

func (a *App) handleAccountsDelete(c *gin.Context) {
	acc, err := a.account(c)
	if errors.Is(err, api.ErrResourceNotFound) {
		c.JSON(http.StatusNotFound, a.error("Account not found"))
		return
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, a.error(err))
		return
	}
	if err := a.api.DeleteAccount(c, acc); err != nil {
		panic(fmt.Errorf("failed to delete account: %w", err))
	}
	c.JSON(http.StatusOK, gin.H{})
}

type AccountAmountMonthInput struct {
	Month string `uri:"month" binding:"yearmonth"`
}

type AccountAmountInput struct {
	Currency string  `json:"currency" binding:"required"`
	Amount   float64 `json:"amount"`
}

func (a *App) handleAccountAmountSet(c *gin.Context) {
	acc, err := a.account(c)
	if errors.Is(err, api.ErrResourceNotFound) {
		c.JSON(http.StatusNotFound, a.error("Account not found"))
		return
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, a.error(err))
		return
	}
	var m AccountAmountMonthInput
	if err := c.ShouldBindUri(&m); err != nil {
		c.JSON(http.StatusBadRequest, a.error(err))
		return
	}
	var input AccountAmountInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, a.error(err))
		return
	}
	if m.Month == "" {
		if err = a.api.SetAccountCurrentAmount(c, acc, input.Currency, input.Amount); err != nil {
			panic(fmt.Errorf("failed to set account amount: %w", err))
		}
		c.JSON(http.StatusOK, gin.H{})
		return
	}
	if err = a.api.SetAccountAmount(c, acc, m.Month, input.Currency, input.Amount); err != nil {
		panic(fmt.Errorf("failed to set account amount: %w", err))
	}
	c.JSON(http.StatusOK, gin.H{})
}

func (a *App) handleAccountAmountGet(c *gin.Context) {
	acc, err := a.account(c)
	if errors.Is(err, api.ErrResourceNotFound) {
		c.JSON(http.StatusNotFound, a.error("Account not found"))
		return
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, a.error(err))
		return
	}
	var m AccountAmountMonthInput
	if err := c.ShouldBindUri(&m); err != nil {
		c.JSON(http.StatusBadRequest, a.error(err))
		return
	}
	if m.Month == "" {
		amts, err := a.api.GetAccountCurrentAmount(c, acc)
		if err != nil {
			panic(fmt.Errorf("failed to get account amount: %w", err))
		}
		c.JSON(http.StatusOK, amts)
		return
	}
	amts, err := a.api.GetAccountAmount(c, acc, m.Month)
	if err != nil {
		panic(fmt.Errorf("failed to get account amount: %w", err))
	}
	c.JSON(http.StatusOK, amts)
}
