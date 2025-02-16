package service

import (
	"context"

	"github.com/esklo/avito-backend-winter-2025/internal/model"
)

//go:generate mockgen -destination=../../mocks/mock_service.go -package=mocks github.com/esklo/avito-backend-winter-2025/internal/service Hasher,Authenticator,UserManager,Shop

type Hasher interface {
	Hash(password string) (hash []byte, salt []byte, err error)
	Verify(password string, hash []byte, salt []byte) bool
}

type Authenticator interface {
	Login(ctx context.Context, username, password string) (string, error)
	ValidateToken(ctx context.Context, token string) (string, error)
}

type UserManager interface {
	Create(ctx context.Context, username, password string) (*model.User, error)
	Info(ctx context.Context, username string) (*model.Info, error)
	Transfer(ctx context.Context, from, to string, amount int) error
}

type Shop interface {
	GetItem(ctx context.Context, name string) (*model.Item, error)
	BuyItem(ctx context.Context, name, username string) error
}
