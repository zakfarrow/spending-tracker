package main

import (
	"net/http"
	"spending-tracker/models"
	"spending-tracker/templates"
	"spending-tracker/templates/components"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Serve static files
	r.Static("/static", "./static")

	// Root endpoint
	r.GET("/", func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		period := models.Period{Month: "March", Year: 2026}
		initial_state := models.AppState{ExpensePeriod: period}
		templates.Index(initial_state).Render(c.Request.Context(), c.Writer)
	})

	// Healthcheck endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	r.GET("current-expense-period", func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		period := models.Period{Month: "January", Year: 2026}
		components.ExpensePeriodHeading(period).Render(c.Request.Context(), c.Writer)
	})

	// r.POST("/expense", func(c *gin.Context) {
	// 	c.Header("Content-Type", "text/html; charset=utf-8")
	// 	var period models.Period
	// 	if err := c.ShouldBind(&period); err != nil {
	// 		c.String(400, "bad request: %v", err)
	// 		return
	// 	}
	// })

	r.Run(":8080")
}
