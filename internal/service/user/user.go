package user

import (
	"context"
	"fmt"

	"github.com/esklo/avito-backend-winter-2025/internal/model"
	"github.com/esklo/avito-backend-winter-2025/internal/repository"
	"github.com/esklo/avito-backend-winter-2025/internal/service"
)

var _ service.UserManager = (*Service)(nil)

type Service struct {
	repo            repository.Repository
	passwordService service.Hasher
}

func NewService(repo repository.Repository, passwordService service.Hasher) *Service {
	return &Service{repo: repo, passwordService: passwordService}
}

func (s *Service) Create(ctx context.Context, username, password string) (user *model.User, err error) {
	if username == "" || password == "" {
		return nil, model.ErrBadRequest
	}

	user = &model.User{
		Username: username,
	}

	user.Password, user.Salt, err = s.passwordService.Hash(password)
	if err != nil {
		return nil, model.ErrInternalServerError
	}

	err = s.repo.CreateUser(ctx, nil, user)
	if err != nil {
		return nil, model.ErrInternalServerError
	}

	return s.repo.FindUser(ctx, nil, user.Username)
}

func (s *Service) Info(ctx context.Context, username string) (*model.Info, error) {
	if username == "" {
		return nil, model.ErrBadRequest
	}

	var info *model.Info

	err := s.repo.WithTx(ctx, func(tx repository.DB) error {
		user, err := s.repo.FindUser(ctx, tx, username)
		if err != nil {
			return model.ErrUnauthorized
		}

		inventory, err := s.repo.ListInventory(ctx, tx, user.ID)
		if err != nil {
			return fmt.Errorf("%w: can not get inventory: %w", model.ErrInternalServerError, err)
		}

		history, err := s.repo.ListTransactions(ctx, tx, user.ID)
		if err != nil {
			return fmt.Errorf("%w: can not get transaction history: %w", model.ErrInternalServerError, err)
		}

		info = &model.Info{
			Coins:       user.Balance,
			Inventory:   inventory,
			CoinHistory: history,
		}

		return nil
	})
	if err != nil {
		return nil, model.ErrInternalServerError
	}

	return info, nil
}

func (s *Service) Transfer(ctx context.Context, from, to string, amount int) error {
	if from == "" {
		return model.ErrUnauthorized
	}

	if amount <= 0 || to == "" || from == to {
		return model.ErrBadRequest
	}

	return s.repo.WithTx(ctx, func(tx repository.DB) error {
		sender, err := s.repo.FindUser(ctx, tx, from)
		if err != nil {
			return model.ErrUnauthorized
		}

		if sender.Balance < amount {
			return model.ErrInsufficientFunds
		}

		receiver, err := s.repo.FindUser(ctx, tx, to)
		if err != nil {
			return model.ErrBadRequest
		}

		return s.repo.MakeTransfer(ctx, tx, sender.ID, receiver.ID, amount)
	})
}
