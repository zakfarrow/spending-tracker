package db

import (
	"context"
	"time"

	"spending-tracker/models"
)

func GetExpensesByPeriod(ctx context.Context, year, month int) ([]models.Expense, error) {
	rows, err := Pool.Query(ctx, `
		SELECT e.id, e.description, e.amount, e.category_id, e.expense_type,
		       e.year, e.month, e.recurring_expense_id, e.created_at, e.updated_at,
		       c.id, c.name, c.color, c.created_at
		FROM expenses e
		LEFT JOIN categories c ON e.category_id = c.id
		WHERE e.year = $1 AND e.month = $2
		ORDER BY e.expense_type DESC, e.created_at
	`, year, month)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expenses []models.Expense
	for rows.Next() {
		var e models.Expense
		var catName, catColor *string
		var catCreatedAt *time.Time
		var cID *int64

		if err := rows.Scan(
			&e.ID, &e.Description, &e.Amount, &e.CategoryID, &e.Type,
			&e.Year, &e.Month, &e.RecurringExpenseID, &e.CreatedAt, &e.UpdatedAt,
			&cID, &catName, &catColor, &catCreatedAt,
		); err != nil {
			return nil, err
		}

		if cID != nil && catName != nil && catColor != nil {
			e.Category = &models.Category{
				ID:    *cID,
				Name:  *catName,
				Color: *catColor,
			}
		}
		expenses = append(expenses, e)
	}
	return expenses, rows.Err()
}

func GetExpensesByPeriodAndType(ctx context.Context, year, month int, expenseType models.ExpenseType) ([]models.Expense, error) {
	rows, err := Pool.Query(ctx, `
		SELECT e.id, e.description, e.amount, e.category_id, e.expense_type,
		       e.year, e.month, e.recurring_expense_id, e.created_at, e.updated_at,
		       c.id, c.name, c.color, c.created_at
		FROM expenses e
		LEFT JOIN categories c ON e.category_id = c.id
		WHERE e.year = $1 AND e.month = $2 AND e.expense_type = $3
		ORDER BY e.created_at
	`, year, month, expenseType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expenses []models.Expense
	for rows.Next() {
		var e models.Expense
		var cID *int64
		var catName, catColor *string
		var catCreatedAt *time.Time

		if err := rows.Scan(
			&e.ID, &e.Description, &e.Amount, &e.CategoryID, &e.Type,
			&e.Year, &e.Month, &e.RecurringExpenseID, &e.CreatedAt, &e.UpdatedAt,
			&cID, &catName, &catColor, &catCreatedAt,
		); err != nil {
			return nil, err
		}

		if cID != nil && catName != nil && catColor != nil {
			e.Category = &models.Category{
				ID:    *cID,
				Name:  *catName,
				Color: *catColor,
			}
		}
		expenses = append(expenses, e)
	}
	return expenses, rows.Err()
}

func GetExpenseByID(ctx context.Context, id int64) (*models.Expense, error) {
	var e models.Expense
	var cID *int64
	var catName, catColor *string
	var catCreatedAt *time.Time

	err := Pool.QueryRow(ctx, `
		SELECT e.id, e.description, e.amount, e.category_id, e.expense_type,
		       e.year, e.month, e.recurring_expense_id, e.created_at, e.updated_at,
		       c.id, c.name, c.color, c.created_at
		FROM expenses e
		LEFT JOIN categories c ON e.category_id = c.id
		WHERE e.id = $1
	`, id).Scan(
		&e.ID, &e.Description, &e.Amount, &e.CategoryID, &e.Type,
		&e.Year, &e.Month, &e.RecurringExpenseID, &e.CreatedAt, &e.UpdatedAt,
		&cID, &catName, &catColor, &catCreatedAt,
	)
	if err != nil {
		return nil, err
	}

	if cID != nil && catName != nil && catColor != nil {
		e.Category = &models.Category{
			ID:    *cID,
			Name:  *catName,
			Color: *catColor,
		}
	}
	return &e, nil
}

func CreateExpense(ctx context.Context, expense models.Expense) (*models.Expense, error) {
	var e models.Expense
	err := Pool.QueryRow(ctx, `
		INSERT INTO expenses (description, amount, category_id, expense_type, year, month, recurring_expense_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, description, amount, category_id, expense_type, year, month, recurring_expense_id, created_at, updated_at
	`, expense.Description, expense.Amount, expense.CategoryID, expense.Type,
		expense.Year, expense.Month, expense.RecurringExpenseID,
	).Scan(&e.ID, &e.Description, &e.Amount, &e.CategoryID, &e.Type,
		&e.Year, &e.Month, &e.RecurringExpenseID, &e.CreatedAt, &e.UpdatedAt)
	if err != nil {
		return nil, err
	}

	if e.CategoryID != nil {
		cat, err := GetCategoryByID(ctx, *e.CategoryID)
		if err == nil {
			e.Category = cat
		}
	}
	return &e, nil
}

func UpdateExpense(ctx context.Context, id int64, description string, amount float64, categoryID *int64, expenseType models.ExpenseType) error {
	_, err := Pool.Exec(ctx, `
		UPDATE expenses
		SET description = $2, amount = $3, category_id = $4, expense_type = $5, updated_at = NOW()
		WHERE id = $1
	`, id, description, amount, categoryID, expenseType)
	return err
}

func DeleteExpense(ctx context.Context, id int64) error {
	_, err := Pool.Exec(ctx, `DELETE FROM expenses WHERE id = $1`, id)
	return err
}
