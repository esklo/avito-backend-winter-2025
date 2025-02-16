package repository

import (
	"context"
	"fmt"

	"github.com/esklo/avito-backend-winter-2025/internal/model"
)

func (r *repo) FindUser(ctx context.Context, tx DB, username string) (*model.User, error) {
	db := r.getExecutor(tx)

	var user model.User

	err := db.QueryRow(ctx, `
		SELECT id, username, password, salt, balance 
		FROM users 
		WHERE username = $1;
	`, username).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Salt,
		&user.Balance,
	)
	if err != nil {
		return nil, fmt.Errorf("select user: %w", err)
	}

	return &user, nil
}

func (r *repo) CreateUser(ctx context.Context, tx DB, user *model.User) error {
	db := r.getExecutor(tx)

	_, err := db.Exec(ctx, `
		INSERT INTO users (username, password, salt) 
		VALUES ($1, $2, $3);
	`, user.Username, user.Password, user.Salt)
	if err != nil {
		return fmt.Errorf("insert user: %w", err)
	}

	return nil
}
