package models

import "time"

type ExpenseType string

const (
	ExpenseTypeRecurring ExpenseType = "recurring"
	ExpenseTypeOneTime   ExpenseType = "one_time"
)

type Expense struct {
	ID                 int64       `json:"id"`
	Description        string      `json:"description"`
	Amount             float64     `json:"amount"`
	CategoryID         *int64      `json:"category_id"`
	Category           *Category   `json:"category,omitempty"`
	Type               ExpenseType `json:"expense_type"`
	Year               int         `json:"year"`
	Month              int         `json:"month"`
	RecurringExpenseID *int64      `json:"recurring_expense_id,omitempty"`
	CreatedAt          time.Time   `json:"created_at"`
	UpdatedAt          time.Time   `json:"updated_at"`
}

func (e Expense) IsRecurring() bool {
	return e.Type == ExpenseTypeRecurring
}
