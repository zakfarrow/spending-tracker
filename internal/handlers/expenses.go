package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"spending-tracker/db"
	"spending-tracker/models"
	"spending-tracker/templates/components"
)

// GetExpenses returns the expense list for a given period and filter
func (h *Handler) GetExpenses(c *gin.Context) {
	year, _ := strconv.Atoi(c.Query("year"))
	month, _ := strconv.Atoi(c.Query("month"))
	filter := models.ExpenseFilter(c.DefaultQuery("filter", "all"))

	period := models.Period{Year: year, Month: month}
	state, err := h.loadAppState(c.Request.Context(), period, filter)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error loading expenses: %v", err)
		return
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	components.ExpenseList(state.Expenses, state.Categories, period, filter).Render(c.Request.Context(), c.Writer)
}

// CreateExpense creates a new expense
func (h *Handler) CreateExpense(c *gin.Context) {
	year, _ := strconv.Atoi(c.PostForm("year"))
	month, _ := strconv.Atoi(c.PostForm("month"))
	description := c.PostForm("description")
	amount, _ := strconv.ParseFloat(c.PostForm("amount"), 64)
	categoryID, _ := strconv.ParseInt(c.PostForm("category_id"), 10, 64)
	expenseType := models.ExpenseType(c.PostForm("expense_type"))

	var catIDPtr *int64
	if categoryID > 0 {
		catIDPtr = &categoryID
	}

	expense := models.Expense{
		Description: description,
		Amount:      amount,
		CategoryID:  catIDPtr,
		Type:        expenseType,
		Year:        year,
		Month:       month,
	}

	created, err := db.CreateExpense(c.Request.Context(), expense)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error creating expense: %v", err)
		return
	}

	if expenseType == models.ExpenseTypeRecurring {
		db.CreateRecurringExpense(c.Request.Context(), description, amount, catIDPtr)
	}

	period := models.Period{Year: year, Month: month}
	state, err := h.loadAppState(c.Request.Context(), period, models.FilterAll)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error loading data: %v", err)
		return
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	components.ExpenseRowWithOOB(*created, state.Categories, state.Summary, state.Income, period).Render(c.Request.Context(), c.Writer)
}

// UpdateExpense updates an existing expense
func (h *Handler) UpdateExpense(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	description := c.PostForm("description")
	amount, _ := strconv.ParseFloat(c.PostForm("amount"), 64)
	categoryID, _ := strconv.ParseInt(c.PostForm("category_id"), 10, 64)
	expenseType := models.ExpenseType(c.PostForm("expense_type"))
	year, _ := strconv.Atoi(c.PostForm("year"))
	month, _ := strconv.Atoi(c.PostForm("month"))

	var catIDPtr *int64
	if categoryID > 0 {
		catIDPtr = &categoryID
	}

	if err := db.UpdateExpense(c.Request.Context(), id, description, amount, catIDPtr, expenseType); err != nil {
		c.String(http.StatusInternalServerError, "Error updating expense: %v", err)
		return
	}

	period := models.Period{Year: year, Month: month}
	state, err := h.loadAppState(c.Request.Context(), period, models.FilterAll)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error loading data: %v", err)
		return
	}

	expense, err := db.GetExpenseByID(c.Request.Context(), id)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error loading expense: %v", err)
		return
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	components.ExpenseRowWithOOB(*expense, state.Categories, state.Summary, state.Income, period).Render(c.Request.Context(), c.Writer)
}

// DeleteExpense deletes an expense
func (h *Handler) DeleteExpense(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	year, _ := strconv.Atoi(c.Query("year"))
	month, _ := strconv.Atoi(c.Query("month"))

	if err := db.DeleteExpense(c.Request.Context(), id); err != nil {
		c.String(http.StatusInternalServerError, "Error deleting expense: %v", err)
		return
	}

	period := models.Period{Year: year, Month: month}
	state, err := h.loadAppState(c.Request.Context(), period, models.FilterAll)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error loading data: %v", err)
		return
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	components.SummaryOOB(state.Summary, state.Income, period).Render(c.Request.Context(), c.Writer)
}

// ExpenseModal returns the add expense modal form
func (h *Handler) ExpenseModal(c *gin.Context) {
	year, _ := strconv.Atoi(c.Query("year"))
	month, _ := strconv.Atoi(c.Query("month"))

	period := models.Period{Year: year, Month: month}
	categories, err := db.GetAllCategories(c.Request.Context())
	if err != nil {
		c.String(http.StatusInternalServerError, "Error loading categories: %v", err)
		return
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	components.AddExpenseModal(categories, period).Render(c.Request.Context(), c.Writer)
}
