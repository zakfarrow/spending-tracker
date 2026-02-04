package main

import (
	"net/http"
	"spending-tracker/templates"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Serve static files
	r.Static("/static", "./static")

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

	// r.POST("/expense", func(c *gin.Context) {
	// 	c.Header("Content-Type", "text/html; charset=utf-8")
	// 	var period Period
	// 	if err := c.ShouldBind(&period); err != nil {
	// 		c.String(400, "bad request: %v", err)
	// 		return
	// 	}
	// })

	r.Run(":8080")
}
