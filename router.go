package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func route(router *gin.Engine) {
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "Hello world"})
	})
}
