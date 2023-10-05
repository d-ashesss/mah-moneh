package rest

import (
	"fmt"
	"github.com/d-ashesss/mah-moneh/internal/categories"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"net/http"
	"time"
)

type CreateCategoryInput struct {
	Name string `json:"name" binding:"required"`
}

func (i *CreateCategoryInput) Bind(c *gin.Context) error {
	return NewErrBadRequestOrNil(c.ShouldBind(i))
}

type GetCategoryInput struct {
	UUID string `uri:"uuid" binding:"required,uuid"`
}

func (i *GetCategoryInput) Bind(c *gin.Context) error {
	return NewErrBadRequestOrNil(c.ShouldBindUri(i))
}

func (h *handler) category(c *gin.Context) (*categories.Category, error) {
	var input GetCategoryInput
	if err := input.Bind(c); err != nil {
		return nil, err
	}
	cat, err := h.categories.GetCategory(c, uuid.FromStringOrNil(input.UUID))
	if err != nil {
		return nil, err
	}
	if cat.User.ID != h.user(c).ID {
		return nil, ErrResourceNotFound
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

func NewListCategoriesResponse(cats []*categories.Category) []*CategoryResponse {
	r := make([]*CategoryResponse, 0, len(cats))
	for _, cat := range cats {
		r = append(r, NewCategoryResponse(cat))
	}
	return r
}

func (h *handler) handleCategoriesCreate(c *gin.Context) {
	var input CreateCategoryInput
	if err := input.Bind(c); err != nil {
		h.handleError(c, err)
		return
	}
	cat, err := h.categories.CreateCategory(c, h.user(c), input.Name)
	if err != nil {
		h.handleError(c, fmt.Errorf("failed to create category: %w", err))
		return
	}
	c.JSON(http.StatusCreated, NewCategoryResponse(cat))
}

func (h *handler) handleCategoriesList(c *gin.Context) {
	cats, err := h.categories.GetUserCategories(c, h.user(c))
	if err != nil {
		h.handleError(c, fmt.Errorf("failed to get user categories: %w", err))
		return
	}
	c.JSON(http.StatusOK, NewListCategoriesResponse(cats))
}

func (h *handler) handleCategoriesDelete(c *gin.Context) {
	cat, err := h.category(c)
	if err != nil {
		h.handleError(c, fmt.Errorf("failed to find category: %w", err))
		return
	}
	if err := h.categories.DeleteCategory(c, cat); err != nil {
		h.handleError(c, fmt.Errorf("failed to delete category: %w", err))
		return
	}
	c.Status(http.StatusNoContent)
}
