package main

import (
	"net/http"

	"spending-tracker/template"
	"spending-tracker/template/component"

	"github.com/gin-gonic/gin"
)

type Period struct {
	Month uint8  `form:"month" binding:"required,gte=1,lte=12"`
	Year  uint16 `form:"year" binding:"required"`
}

func main() {
	r := gin.Default()

	// Root endpoint
	r.GET("/", func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		template.Index().Render(c.Request.Context(), c.Writer)
	})

	// Healthcheck endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	r.POST("/expense", func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		var period Period
		if err := c.ShouldBind(&period); err != nil {
			c.String(400, "bad request: %v", err)
			return
		}

		component.ExpenseDetail(period.Month, period.Year).Render(c.Request.Context(), c.Writer)
	})

	r.Run(":8080")
}
