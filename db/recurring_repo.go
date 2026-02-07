package db

import (
	"context"
	"spending-tracker/models"
)

func GetActiveRecurringExpenses(ctx context.Context) ([]models.Expense, error) {
	rows, err := Pool.Query(ctx, `
		SELECT r.id, r.description, r.amount, r.category_id,
		       c.id, c.name, c.color
		FROM recurring_expenses r
		LEFT JOIN categories c ON r.category_id = c.id
		WHERE r.is_active = true
		ORDER BY r.created_at
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expenses []models.Expense
	for rows.Next() {
		var e models.Expense
		var cID *int64
		var catName, catColor *string

		if err := rows.Scan(
			&e.RecurringExpenseID, &e.Description, &e.Amount, &e.CategoryID,
			&cID, &catName, &catColor,
		); err != nil {
			return nil, err
		}

		e.Type = models.ExpenseTypeRecurring
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

func IsMonthInitialized(ctx context.Context, year, month int) (bool, error) {
	var exists bool
	err := Pool.QueryRow(ctx, `
		SELECT EXISTS(SELECT 1 FROM initialized_months WHERE year = $1 AND month = $2)
	`, year, month).Scan(&exists)
	return exists, err
}

func MarkMonthInitialized(ctx context.Context, year, month int) error {
	_, err := Pool.Exec(ctx, `
		INSERT INTO initialized_months (year, month)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING
	`, year, month)
	return err
}

func InitializeMonth(ctx context.Context, year, month int) error {
	initialized, err := IsMonthInitialized(ctx, year, month)
	if err != nil {
		return err
	}
	if initialized {
		return nil
	}

	recurring, err := GetActiveRecurringExpenses(ctx)
	if err != nil {
		return err
	}

	for _, r := range recurring {
		expense := models.Expense{
			Description:        r.Description,
			Amount:             r.Amount,
			CategoryID:         r.CategoryID,
			Type:               models.ExpenseTypeRecurring,
			Year:               year,
			Month:              month,
			RecurringExpenseID: r.RecurringExpenseID,
		}
		_, err := CreateExpense(ctx, expense)
		if err != nil {
			return err
		}
	}

	prevPeriod := models.Period{Year: year, Month: month}.Prev()
	prevIncome, _ := GetIncomeByPeriod(ctx, prevPeriod.Year, prevPeriod.Month)
	if prevIncome > 0 {
		UpsertIncome(ctx, year, month, prevIncome)
	}

	return MarkMonthInitialized(ctx, year, month)
}

func CreateRecurringExpense(ctx context.Context, description string, amount float64, categoryID *int64) (int64, error) {
	var id int64
	err := Pool.QueryRow(ctx, `
		INSERT INTO recurring_expenses (description, amount, category_id)
		VALUES ($1, $2, $3)
		RETURNING id
	`, description, amount, categoryID).Scan(&id)
	return id, err
}

func DeleteRecurringExpense(ctx context.Context, id int64) error {
	_, err := Pool.Exec(ctx, `
		UPDATE recurring_expenses SET is_active = false WHERE id = $1
	`, id)
	return err
}
