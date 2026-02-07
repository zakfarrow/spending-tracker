package db

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func GetIncomeByPeriod(ctx context.Context, year, month int) (float64, error) {
	var amount float64
	err := Pool.QueryRow(ctx, `
		SELECT amount
		FROM income
		WHERE year = $1 AND month = $2
	`, year, month).Scan(&amount)
	if err == pgx.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return amount, nil
}

func UpsertIncome(ctx context.Context, year, month int, amount float64) error {
	_, err := Pool.Exec(ctx, `
		INSERT INTO income (year, month, amount)
		VALUES ($1, $2, $3)
		ON CONFLICT (year, month)
		DO UPDATE SET amount = $3, updated_at = NOW()
	`, year, month, amount)
	return err
}
