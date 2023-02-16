package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (a *App) registerHttpHandlers() {
	r := gin.Default()

	r.GET("/", a.handleIndex)

	a.server.Handler = r
}

func (a *App) handleIndex(c *gin.Context) {
	c.String(http.StatusOK, "ok")
}
