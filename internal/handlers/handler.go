package handlers

import (
	"context"
	"spending-tracker/db"
	"spending-tracker/models"
)

// Handler handles HTTP requests for the application
type Handler struct{}

// NewHandler creates a new Handler instance
func NewHandler() *Handler {
	return &Handler{}
}

// loadAppState loads the complete application state for a given period and filter
func (h *Handler) loadAppState(ctx context.Context, period models.Period, filter models.ExpenseFilter) (models.AppState, error) {
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
