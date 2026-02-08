package main

import (
	"context"
	"log"
	"spending-tracker/db"
	"spending-tracker/internal/handlers"

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

	// Initialize handler
	h := handlers.NewHandler()

	// Page routes
	r.GET("/", h.Index)
	r.GET("/period/:year/:month", h.Period)
	r.GET("/health", h.Health)

	// Income routes
	r.PUT("/income", h.UpdateIncome)

	// Expense routes
	r.GET("/expenses", h.GetExpenses)
	r.POST("/expenses", h.CreateExpense)
	r.PUT("/expenses/:id", h.UpdateExpense)
	r.DELETE("/expenses/:id", h.DeleteExpense)

	// Category routes
	r.GET("/categories", h.GetCategories)
	r.GET("/categories/options", h.GetCategoryOptions)
	r.POST("/categories", h.CreateCategory)
	r.PUT("/categories/:id", h.UpdateCategory)
	r.DELETE("/categories/:id", h.DeleteCategory)

	// Modal routes
	r.GET("/modals/expense", h.ExpenseModal)
	r.GET("/modals/category", h.CategoryModal)

	r.Run(":8080")
}
