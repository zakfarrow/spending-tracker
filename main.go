package main

import (
	"context"
	"log"
	"net/http"
	"spending-tracker/db"
	"spending-tracker/models"
	"spending-tracker/templates"
	"spending-tracker/templates/components"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	ctx := context.Background()

	if err := db.Connect(ctx); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.RunMigrations(ctx); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	r := gin.Default()
	r.Static("/static", "./static")

	r.GET("/", handleIndex)
	r.GET("/period/:year/:month", handlePeriod)
	r.GET("/health", handleHealth)

	r.PUT("/income", handleUpdateIncome)

	r.GET("/expenses", handleGetExpenses)
	r.POST("/expenses", handleCreateExpense)
	r.PUT("/expenses/:id", handleUpdateExpense)
	r.DELETE("/expenses/:id", handleDeleteExpense)

	r.GET("/categories", handleGetCategories)
	r.POST("/categories", handleCreateCategory)
	r.PUT("/categories/:id", handleUpdateCategory)
	r.DELETE("/categories/:id", handleDeleteCategory)

	r.Run(":8080")
}

func loadAppState(ctx context.Context, period models.Period, filter models.ExpenseFilter) (models.AppState, error) {
	if err := db.InitializeMonth(ctx, period.Year, period.Month); err != nil {
		return models.AppState{}, err
	}

	income, err := db.GetIncomeByPeriod(ctx, period.Year, period.Month)
	if err != nil {
		return models.AppState{}, err
	}

	var expenses []models.Expense
	switch filter {
	case models.FilterRecurring:
		expenses, err = db.GetExpensesByPeriodAndType(ctx, period.Year, period.Month, models.ExpenseTypeRecurring)
	case models.FilterOneTime:
		expenses, err = db.GetExpensesByPeriodAndType(ctx, period.Year, period.Month, models.ExpenseTypeOneTime)
	default:
		expenses, err = db.GetExpensesByPeriod(ctx, period.Year, period.Month)
	}
	if err != nil {
		return models.AppState{}, err
	}

	allExpenses, err := db.GetExpensesByPeriod(ctx, period.Year, period.Month)
	if err != nil {
		return models.AppState{}, err
	}

	categories, err := db.GetAllCategories(ctx)
	if err != nil {
		return models.AppState{}, err
	}

	summary := models.CalculateSummary(income, allExpenses, period.DaysInMonth())

	return models.AppState{
		Period:     period,
		Income:     income,
		Expenses:   expenses,
		Categories: categories,
		Summary:    summary,
		Filter:     filter,
	}, nil
}

func handleIndex(c *gin.Context) {
	period := models.CurrentPeriod()
	state, err := loadAppState(c.Request.Context(), period, models.FilterAll)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error loading data: %v", err)
		return
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	templates.Index(state).Render(c.Request.Context(), c.Writer)
}

func handlePeriod(c *gin.Context) {
	year, _ := strconv.Atoi(c.Param("year"))
	month, _ := strconv.Atoi(c.Param("month"))
	filter := models.ExpenseFilter(c.DefaultQuery("filter", "all"))

	period := models.Period{Year: year, Month: month}
	state, err := loadAppState(c.Request.Context(), period, filter)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error loading data: %v", err)
		return
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	components.FullPageContent(state).Render(c.Request.Context(), c.Writer)
}

func handleHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func handleUpdateIncome(c *gin.Context) {
	year, _ := strconv.Atoi(c.PostForm("year"))
	month, _ := strconv.Atoi(c.PostForm("month"))
	amount, _ := strconv.ParseFloat(c.PostForm("amount"), 64)

	if err := db.UpsertIncome(c.Request.Context(), year, month, amount); err != nil {
		c.String(http.StatusInternalServerError, "Error updating income: %v", err)
		return
	}

	period := models.Period{Year: year, Month: month}
	state, err := loadAppState(c.Request.Context(), period, models.FilterAll)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error loading data: %v", err)
		return
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	components.IncomeWithOOB(state).Render(c.Request.Context(), c.Writer)
}

func handleGetExpenses(c *gin.Context) {
	year, _ := strconv.Atoi(c.Query("year"))
	month, _ := strconv.Atoi(c.Query("month"))
	filter := models.ExpenseFilter(c.DefaultQuery("filter", "all"))

	period := models.Period{Year: year, Month: month}
	state, err := loadAppState(c.Request.Context(), period, filter)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error loading expenses: %v", err)
		return
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	components.ExpenseList(state.Expenses, state.Categories, period, filter).Render(c.Request.Context(), c.Writer)
}

func handleCreateExpense(c *gin.Context) {
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
	state, err := loadAppState(c.Request.Context(), period, models.FilterAll)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error loading data: %v", err)
		return
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	components.ExpenseRowWithOOB(*created, state.Categories, state.Summary, state.Income, period).Render(c.Request.Context(), c.Writer)
}

func handleUpdateExpense(c *gin.Context) {
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
	state, err := loadAppState(c.Request.Context(), period, models.FilterAll)
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

func handleDeleteExpense(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	year, _ := strconv.Atoi(c.Query("year"))
	month, _ := strconv.Atoi(c.Query("month"))

	if err := db.DeleteExpense(c.Request.Context(), id); err != nil {
		c.String(http.StatusInternalServerError, "Error deleting expense: %v", err)
		return
	}

	period := models.Period{Year: year, Month: month}
	state, err := loadAppState(c.Request.Context(), period, models.FilterAll)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error loading data: %v", err)
		return
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	components.SummaryOOB(state.Summary, state.Income, period).Render(c.Request.Context(), c.Writer)
}

func handleGetCategories(c *gin.Context) {
	categories, err := db.GetAllCategories(c.Request.Context())
	if err != nil {
		c.String(http.StatusInternalServerError, "Error loading categories: %v", err)
		return
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	components.CategoryList(categories).Render(c.Request.Context(), c.Writer)
}

func handleCreateCategory(c *gin.Context) {
	name := c.PostForm("name")
	color := c.PostForm("color")

	_, err := db.CreateCategory(c.Request.Context(), name, color)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error creating category: %v", err)
		return
	}

	categories, err := db.GetAllCategories(c.Request.Context())
	if err != nil {
		c.String(http.StatusInternalServerError, "Error loading categories: %v", err)
		return
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	components.CategoryList(categories).Render(c.Request.Context(), c.Writer)
}

func handleUpdateCategory(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	name := c.PostForm("name")
	color := c.PostForm("color")

	if err := db.UpdateCategory(c.Request.Context(), id, name, color); err != nil {
		c.String(http.StatusInternalServerError, "Error updating category: %v", err)
		return
	}

	categories, err := db.GetAllCategories(c.Request.Context())
	if err != nil {
		c.String(http.StatusInternalServerError, "Error loading categories: %v", err)
		return
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	components.CategoryList(categories).Render(c.Request.Context(), c.Writer)
}

func handleDeleteCategory(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	if err := db.DeleteCategory(c.Request.Context(), id); err != nil {
		c.String(http.StatusInternalServerError, "Error deleting category: %v", err)
		return
	}

	categories, err := db.GetAllCategories(c.Request.Context())
	if err != nil {
		c.String(http.StatusInternalServerError, "Error loading categories: %v", err)
		return
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	components.CategoryList(categories).Render(c.Request.Context(), c.Writer)
}
