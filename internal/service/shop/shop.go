package shop

import (
	"context"
	"fmt"

	"github.com/esklo/avito-backend-winter-2025/internal/model"
	"github.com/esklo/avito-backend-winter-2025/internal/repository"
	"github.com/esklo/avito-backend-winter-2025/internal/service"
)

var _ service.Shop = (*Service)(nil)

type Service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetItem(ctx context.Context, name string) (*model.Item, error) {
	return s.repo.FindItem(ctx, nil, name)
}

func (s *Service) BuyItem(ctx context.Context, name, username string) error {
	if username == "" {
		return model.ErrUnauthorized
	}

	if name == "" {
		return model.ErrBadRequest
	}

	return s.repo.WithTx(ctx, func(tx repository.DB) error {
		user, err := s.repo.FindUser(ctx, tx, username)
		if err != nil {
			return model.ErrUnauthorized
		}

		item, err := s.repo.FindItem(ctx, tx, name)
		if err != nil {
			return model.ErrNotFound
		}

		if user.Balance < item.Price {
			return fmt.Errorf("%w: need %d coins, has %d",
				model.ErrInsufficientFunds,
				item.Price,
				user.Balance,
			)
		}

		return s.repo.MakePurchase(ctx, tx, user.ID, item.ID, item.Price)
	})
}
