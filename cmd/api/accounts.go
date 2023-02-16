package main

import (
	"errors"
	"github.com/d-ashesss/mah-moneh/internal/accounts"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"log"
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
		return nil, errors.New(http.StatusText(http.StatusForbidden))
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
		log.Printf("Failed to create account: %s", err)
		c.JSON(http.StatusInternalServerError, a.error(err))
		return
	}

	c.JSON(http.StatusCreated, NewAccountResponse(acc))
}

func (a *App) handleAccountsGet(c *gin.Context) {
	accs, err := a.api.GetUserAccounts(c, a.user(c))
	if err != nil {
		log.Printf("Failed to get user accounts: %s", err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusOK, MapAccountsResponse(accs))
}

func (a *App) handleAccountsUpdate(c *gin.Context) {
	acc, err := a.account(c)
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
		log.Printf("Failed to update account: %s", err)
		c.JSON(http.StatusInternalServerError, a.error(err))
		return
	}
	c.JSON(http.StatusOK, NewAccountResponse(acc))
}

func (a *App) handleAccountsDelete(c *gin.Context) {
	acc, err := a.account(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, a.error(err))
		return
	}
	if err := a.api.DeleteAccount(c, acc); err != nil {
		log.Printf("Failed to delete account: %s", err)
		c.JSON(http.StatusInternalServerError, a.error(err))
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
