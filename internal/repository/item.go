package repository

import (
	"context"
	"fmt"

	"github.com/esklo/avito-backend-winter-2025/internal/model"
)

func (r *repo) FindItem(ctx context.Context, tx DB, name string) (*model.Item, error) {
	db := r.getExecutor(tx)

	var item model.Item

	err := db.QueryRow(ctx, `
		SELECT id, name, price
		FROM items
		WHERE name = $1;
	`, name).Scan(
		&item.ID,
		&item.Name,
		&item.Price,
	)
	if err != nil {
		return nil, fmt.Errorf("select item: %w", err)
	}

	return &item, nil
}
