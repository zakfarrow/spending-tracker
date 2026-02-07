package models

type ExpenseFilter string

const (
	FilterAll       ExpenseFilter = "all"
	FilterRecurring ExpenseFilter = "recurring"
	FilterOneTime   ExpenseFilter = "one_time"
)

type AppState struct {
	Period     Period
	Income     float64
	Expenses   []Expense
	Categories []Category
	Summary    Summary
	Filter     ExpenseFilter
}
