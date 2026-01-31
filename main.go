package main

import (
	"net/http"

	"spending-tracker/templates"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Root endpoint
	r.GET("/", func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		templates.Index().Render(c.Request.Context(), c.Writer)
	})

	// Healthcheck endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	r.Run(":8080")
}
