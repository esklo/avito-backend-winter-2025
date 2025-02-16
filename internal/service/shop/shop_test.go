package shop

import (
	"context"
	"errors"
	"testing"

	"github.com/esklo/avito-backend-winter-2025/internal/model"
	"github.com/esklo/avito-backend-winter-2025/internal/repository"
	"github.com/esklo/avito-backend-winter-2025/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type testSuite struct {
	shop *Service
	repo *mocks.MockRepository
}

func newTestSuite(t *testing.T) *testSuite {
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockRepository(ctrl)

	return &testSuite{
		repo: repo,
		shop: NewService(repo),
	}
}

func TestService_GetItem(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("item found", func(t *testing.T) {
		t.Parallel()
		ts := newTestSuite(t)

		expected := &model.Item{
			ID:    1,
			Name:  "hoody",
			Price: 100,
		}

		ts.repo.EXPECT().
			FindItem(gomock.Any(), nil, "hoody").
			Return(expected, nil)

		item, err := ts.shop.GetItem(ctx, "hoody")
		assert.NoError(t, err)
		assert.Equal(t, expected, item)
	})

	t.Run("item not found", func(t *testing.T) {
		t.Parallel()
		ts := newTestSuite(t)

		ts.repo.EXPECT().
			FindItem(gomock.Any(), nil, "unknown").
			Return(nil, errors.New("not found"))

		item, err := ts.shop.GetItem(ctx, "unknown")
		assert.Error(t, err)
		assert.Nil(t, item)
	})
}

func TestService_BuyItem(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("validation cases", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name          string
			item          string
			username      string
			expectedError error
		}{
			{
				name:          "empty username",
				item:          "item",
				username:      "",
				expectedError: model.ErrUnauthorized,
			},
			{
				name:          "empty item",
				item:          "",
				username:      "user",
				expectedError: model.ErrBadRequest,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				ts := newTestSuite(t)
				err := ts.shop.BuyItem(ctx, tt.item, tt.username)
				assert.ErrorIs(t, err, tt.expectedError)
			})
		}
	})

	t.Run("successful purchase", func(t *testing.T) {
		t.Parallel()
		ts := newTestSuite(t)

		user := &model.User{ID: 1, Username: "buyer", Balance: 1000}
		item := &model.Item{ID: 1, Name: "hoody", Price: 500}

		ts.repo.EXPECT().
			WithTx(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, fn func(repository.DB) error) error {
				return fn(nil)
			})

		ts.repo.EXPECT().
			FindUser(gomock.Any(), nil, user.Username).
			Return(user, nil)

		ts.repo.EXPECT().
			FindItem(gomock.Any(), nil, item.Name).
			Return(item, nil)

		ts.repo.EXPECT().
			MakePurchase(gomock.Any(), nil, user.ID, item.ID, item.Price).
			Return(nil)

		err := ts.shop.BuyItem(ctx, item.Name, user.Username)
		assert.NoError(t, err)
	})

	t.Run("insufficient funds", func(t *testing.T) {
		t.Parallel()
		ts := newTestSuite(t)

		user := &model.User{ID: 1, Username: "buyer", Balance: 100}
		item := &model.Item{ID: 1, Name: "hoody", Price: 500}

		ts.repo.EXPECT().
			WithTx(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, fn func(repository.DB) error) error {
				return fn(nil)
			})

		ts.repo.EXPECT().
			FindUser(gomock.Any(), nil, user.Username).
			Return(user, nil)

		ts.repo.EXPECT().
			FindItem(gomock.Any(), nil, item.Name).
			Return(item, nil)

		err := ts.shop.BuyItem(ctx, item.Name, user.Username)
		assert.ErrorIs(t, err, model.ErrInsufficientFunds)
	})
}
