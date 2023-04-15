package rest

import (
	"errors"
	"fmt"
	"github.com/d-ashesss/mah-moneh/internal/datastore"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
)

var (
	ErrResourceNotFound = datastore.ErrRecordNotFound
)

type BadRequestError struct {
	msg string
	err error
}

func NewErrBadRequest(err error) BadRequestError {
	return BadRequestError{
		msg: "Invalid request input",
		err: err,
	}
}

func NewErrBadRequestOrNil(err error) error {
	if err == nil {
		return nil
	}
	return NewErrBadRequest(err)
}

func (err BadRequestError) Error() string {
	return err.msg
}

func (err BadRequestError) Unwrap() error {
	return err.err
}

func NewErrorResponse(err any) gin.H {
	errorstr := "Internal Error"
	switch v := err.(type) {
	case error:
		errorstr = v.Error()
	case string:
		errorstr = v
	}
	return gin.H{"error": errorstr}
}

func (h *handler) handleError(c *gin.Context, err error) {
	if errors.Is(err, ErrResourceNotFound) {
		c.JSON(http.StatusNotFound, NewErrorResponse("Not found"))
		return
	}
	var validationErr validator.ValidationErrors
	if errors.As(err, &validationErr) {
		c.JSON(http.StatusBadRequest, NewErrorResponse(fmt.Sprintf("Invalid value of '%s'", validationErr[0].Field())))
		return
	}
	var requestErr BadRequestError
	if errors.As(err, &requestErr) {
		c.JSON(http.StatusBadRequest, NewErrorResponse(requestErr))
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewErrorResponse(nil))
		log.Printf("[APP] Unexpected error: %v", err)
		return
	}
}

func (h *handler) handleRecovery(c *gin.Context, err any) {
	c.JSON(http.StatusInternalServerError, NewErrorResponse(nil))
	log.Printf("[APP] Panic recovered: %v", err)
}
