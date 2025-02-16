package repository

import (
	"context"
	"fmt"

	"github.com/esklo/avito-backend-winter-2025/internal/model"
)

func (r *repo) MakePurchase(ctx context.Context, tx DB, userID, itemID, price int) error {
	db := r.getExecutor(tx)

	_, err := db.Exec(ctx, `
		INSERT INTO purchases (user_id, item_id, quantity)
		VALUES ($1, $2, 1)
		ON CONFLICT (user_id, item_id)
		DO UPDATE SET quantity = purchases.quantity + 1;
	`, userID, itemID)
	if err != nil {
		return fmt.Errorf("insert purchase: %w", err)
	}

	_, err = db.Exec(ctx, `
		UPDATE users 
		SET balance = balance - $2
		WHERE id = $1;
	`, userID, price)
	if err != nil {
		return fmt.Errorf("update balance: %w", err)
	}

	return nil
}

func (r *repo) ListInventory(ctx context.Context, tx DB, userID int) (inventory []model.Inventory, err error) {
	db := r.getExecutor(tx)

	rows, err := db.Query(ctx, `
		SELECT items.name, purchases.quantity
		FROM purchases
		JOIN items ON items.id = purchases.item_id
		WHERE purchases.user_id = $1;
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("select inventory: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item model.Inventory

		err := rows.Scan(
			&item.Type,
			&item.Quantity,
		)
		if err != nil {
			return nil, err
		}

		inventory = append(inventory, item)
	}

	return inventory, nil
}
