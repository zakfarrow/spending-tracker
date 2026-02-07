package db

import (
	"context"
	"spending-tracker/models"
)

func GetAllCategories(ctx context.Context) ([]models.Category, error) {
	rows, err := Pool.Query(ctx, `
		SELECT id, name, color, created_at
		FROM categories
		ORDER BY name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Color, &c.CreatedAt); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, rows.Err()
}

func GetCategoryByID(ctx context.Context, id int64) (*models.Category, error) {
	var c models.Category
	err := Pool.QueryRow(ctx, `
		SELECT id, name, color, created_at
		FROM categories
		WHERE id = $1
	`, id).Scan(&c.ID, &c.Name, &c.Color, &c.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func CreateCategory(ctx context.Context, name, color string) (*models.Category, error) {
	var c models.Category
	err := Pool.QueryRow(ctx, `
		INSERT INTO categories (name, color)
		VALUES ($1, $2)
		RETURNING id, name, color, created_at
	`, name, color).Scan(&c.ID, &c.Name, &c.Color, &c.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func UpdateCategory(ctx context.Context, id int64, name, color string) error {
	_, err := Pool.Exec(ctx, `
		UPDATE categories
		SET name = $2, color = $3
		WHERE id = $1
	`, id, name, color)
	return err
}

func DeleteCategory(ctx context.Context, id int64) error {
	_, err := Pool.Exec(ctx, `DELETE FROM categories WHERE id = $1`, id)
	return err
}
