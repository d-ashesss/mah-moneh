package main

import (
	"errors"
	"github.com/d-ashesss/mah-moneh/internal/api"
	"github.com/d-ashesss/mah-moneh/internal/categories"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"net/http"
	"time"
)

type CreateCategoryInput struct {
	Name string `json:"name" binding:"required"`
}

type GetCategoryInput struct {
	UUID string `uri:"uuid" binding:"required,uuid"`
}

func (a *App) category(c *gin.Context) (*categories.Category, error) {
	var input GetCategoryInput
	if err := c.ShouldBindUri(&input); err != nil {
		return nil, err
	}
	cat, err := a.api.GetCategory(c, uuid.FromStringOrNil(input.UUID))
	if err != nil {
		return nil, err
	}
	if cat.User.UUID != a.user(c).UUID {
		return nil, api.ErrResourceNotFound
	}
	return cat, err
}

type CategoryResponse struct {
	UUID      string `json:"uuid"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
}

func NewCategoryResponse(cat *categories.Category) *CategoryResponse {
	return &CategoryResponse{
		UUID:      cat.UUID.String(),
		Name:      cat.Name,
		CreatedAt: cat.CreatedAt.Format(time.DateTime),
	}
}

func MapCategoriesResponse(cats []*categories.Category) []*CategoryResponse {
	r := make([]*CategoryResponse, 0, len(cats))
	for _, cat := range cats {
		r = append(r, NewCategoryResponse(cat))
	}
	return r
}

func (a *App) handleCategoriesCreate(c *gin.Context) {
	var input CreateCategoryInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, a.error(err))
		return
	}
	cat, err := a.api.CreateCategory(c, a.user(c), input.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, a.error(err))
		return
	}
	c.JSON(http.StatusCreated, NewCategoryResponse(cat))
}

func (a *App) handleCategoriesGet(c *gin.Context) {
	cats, err := a.api.GetUserCategories(c, a.user(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, a.error(err))
		return
	}
	c.JSON(http.StatusOK, MapCategoriesResponse(cats))
}

func (a *App) handleCategoriesDelete(c *gin.Context) {
	cat, err := a.category(c)
	if errors.Is(err, api.ErrResourceNotFound) {
		c.JSON(http.StatusNotFound, a.error("Category not found"))
		return
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, a.error(err))
		return
	}
	if err := a.api.DeleteCategory(c, cat); err != nil {
		c.JSON(http.StatusInternalServerError, a.error(err))
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
