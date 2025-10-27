package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// Router

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		
		c.JSON(http.StatusOK, "Welcome to my API")
	})
	router.Static("/attachment", "./attachment")
	router.Run("localhost:8080")
}