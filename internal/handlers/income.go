package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"spending-tracker/db"
	"spending-tracker/models"
	"spending-tracker/templates/components"
)

// UpdateIncome updates the income for a given period
func (h *Handler) UpdateIncome(c *gin.Context) {
	year, _ := strconv.Atoi(c.PostForm("year"))
	month, _ := strconv.Atoi(c.PostForm("month"))
	amount, _ := strconv.ParseFloat(c.PostForm("amount"), 64)

	if err := db.UpsertIncome(c.Request.Context(), year, month, amount); err != nil {
		c.String(http.StatusInternalServerError, "Error updating income: %v", err)
		return
	}

	period := models.Period{Year: year, Month: month}
	state, err := h.loadAppState(c.Request.Context(), period, models.FilterAll)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error loading data: %v", err)
		return
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	components.IncomeWithOOB(state).Render(c.Request.Context(), c.Writer)
}
