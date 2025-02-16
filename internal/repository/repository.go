package repository

import (
	"context"
	"fmt"

	"github.com/esklo/avito-backend-winter-2025/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:generate mockgen -destination=../../mocks/mock_db.go -package=mocks github.com/esklo/avito-backend-winter-2025/internal/repository DB
type DB interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, arguments ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, arguments ...any) pgx.Row
}

//go:generate mockgen -destination=../../mocks/mock_repository.go -package=mocks github.com/esklo/avito-backend-winter-2025/internal/repository Repository
type Repository interface {
	FindUser(ctx context.Context, tx DB, username string) (*model.User, error)
	CreateUser(ctx context.Context, tx DB, user *model.User) error
	MakeTransfer(ctx context.Context, tx DB, senderID, receiverID int, amount int) error
	MakePurchase(ctx context.Context, tx DB, userID, itemID, price int) error

	FindItem(ctx context.Context, tx DB, name string) (*model.Item, error)

	ListInventory(ctx context.Context, tx DB, userID int) ([]model.Inventory, error)
	ListTransactions(ctx context.Context, tx DB, userID int) (*model.CoinHistory, error)

	WithTx(ctx context.Context, fn func(DB) error) error
}

type repo struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) Repository {
	return &repo{db: db}
}

func (r *repo) WithTx(ctx context.Context, fn func(DB) error) error {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	})
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)

			return
		}

		err = tx.Commit(ctx)
	}()

	return fn(tx)
}

func (r *repo) getExecutor(tx DB) DB {
	if tx != nil {
		return tx
	}

	return r.db
}
