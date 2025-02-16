package app

import (
	"context"
	"fmt"

	"github.com/esklo/avito-backend-winter-2025/internal/config"
	"github.com/esklo/avito-backend-winter-2025/internal/di"
	"github.com/esklo/avito-backend-winter-2025/internal/http"
	"github.com/esklo/avito-backend-winter-2025/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	cfg       *config.Config
	db        *pgxpool.Pool
	container *di.Container
}

func New(ctx context.Context) (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	db, err := initDB(ctx, cfg.DB)
	if err != nil {
		return nil, fmt.Errorf("init database: %w", err)
	}

	repo := repository.New(db)

	container := di.New(cfg, repo)

	return &App{
		cfg:       cfg,
		db:        db,
		container: container,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	server := http.NewServer(a.container)

	return server.Run(ctx)
}

func (a *App) Shutdown() error {
	a.db.Close()

	return nil
}

func initDB(ctx context.Context, cfg config.DBConfig) (*pgxpool.Pool, error) {
	db, err := pgxpool.New(ctx, cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("connect to database: %w", err)
	}

	if err := db.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return db, nil
}
