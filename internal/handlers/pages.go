package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"spending-tracker/models"
	"spending-tracker/templates"
	"spending-tracker/templates/components"
)

// Index renders the main index page for the current period
func (h *Handler) Index(c *gin.Context) {
	period := models.CurrentPeriod()
	state, err := h.loadAppState(c.Request.Context(), period, models.FilterAll)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error loading data: %v", err)
		return
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	templates.Index(state).Render(c.Request.Context(), c.Writer)
}

// Period renders the page content for a specific period
func (h *Handler) Period(c *gin.Context) {
	year, _ := strconv.Atoi(c.Param("year"))
	month, _ := strconv.Atoi(c.Param("month"))
	filter := models.ExpenseFilter(c.DefaultQuery("filter", "all"))

	period := models.Period{Year: year, Month: month}
	state, err := h.loadAppState(c.Request.Context(), period, filter)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error loading data: %v", err)
		return
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	components.FullPageContent(state).Render(c.Request.Context(), c.Writer)
}
