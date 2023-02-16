package main

import (
	"github.com/d-ashesss/mah-moneh/internal/users"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"net/http"
	"strings"
)

func (a *App) registerHttpHandlers() {
	r := gin.Default()

	r.GET("/", a.handleIndex)

	r.Use(a.authenticate)
	r.POST("/accounts", a.handleAccountsCreate)
	r.GET("/accounts", a.handleAccountsGet)
	r.PUT("/accounts/:uuid", a.handleAccountsUpdate)
	r.DELETE("/accounts/:uuid", a.handleAccountsDelete)

	a.server.Handler = r
}

func (a *App) handleIndex(c *gin.Context) {
	c.String(http.StatusOK, "ok")
}

func (a *App) authenticate(c *gin.Context) {
	auth := c.GetHeader("Authorization")
	uid, _ := strings.CutPrefix(auth, "Bearer ")
	if uid == "" {
		c.JSON(http.StatusUnauthorized, a.error(http.StatusText(http.StatusUnauthorized)))
		c.Abort()
		return
	}
	user := &users.User{UUID: uuid.FromStringOrNil(uid)}
	c.Set("user", user)
	c.Next()
}

func (a *App) error(err any) gin.H {
	errorstr := "Internal Error"
	switch v := err.(type) {
	case error:
		errorstr = v.Error()
	case string:
		errorstr = v
	}
	return gin.H{"error": errorstr}
}

func (a *App) user(c *gin.Context) *users.User {
	user, ok := c.Get("user")
	if !ok {
		return nil
	}
	return user.(*users.User)
}
