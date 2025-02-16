package user

import (
	"context"
	"testing"

	"github.com/esklo/avito-backend-winter-2025/internal/model"
	"github.com/esklo/avito-backend-winter-2025/internal/repository"
	"github.com/esklo/avito-backend-winter-2025/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type testSuite struct {
	users  *Service
	repo   *mocks.MockRepository
	hasher *mocks.MockHasher
}

func newTestSuite(t *testing.T) *testSuite {
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockRepository(ctrl)
	hasher := mocks.NewMockHasher(ctrl)

	return &testSuite{
		repo:   repo,
		hasher: hasher,
		users:  NewService(repo, hasher),
	}
}

func TestService_Create(t *testing.T) {
	t.Parallel()

	ts := newTestSuite(t)
	ctx := context.Background()

	t.Run("validation errors", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name     string
			username string
			password string
		}{
			{"empty both", "", ""},
			{"empty password", "user", ""},
			{"empty username", "", "pass"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				user, err := ts.users.Create(ctx, tt.username, tt.password)
				assert.ErrorIs(t, err, model.ErrBadRequest)
				assert.Nil(t, user)
			})
		}
	})

	t.Run("successful creation", func(t *testing.T) {
		t.Parallel()
		expected := &model.User{
			ID:       1,
			Username: "user",
			Password: []byte("hashed"),
			Salt:     []byte("salt"),
		}

		ts.hasher.EXPECT().
			Hash(gomock.Any()).
			Return(expected.Password, expected.Salt, nil)

		ts.repo.EXPECT().
			CreateUser(gomock.Any(), nil, gomock.Any()).
			Return(nil)

		ts.repo.EXPECT().
			FindUser(gomock.Any(), nil, expected.Username).
			Return(expected, nil)

		user, err := ts.users.Create(ctx, expected.Username, "password")
		assert.NoError(t, err)
		assert.Equal(t, expected, user)
	})
}

func TestService_Info(t *testing.T) {
	t.Parallel()

	ts := newTestSuite(t)
	ctx := context.Background()

	t.Run("empty username", func(t *testing.T) {
		t.Parallel()
		info, err := ts.users.Info(ctx, "")
		assert.ErrorIs(t, err, model.ErrBadRequest)
		assert.Nil(t, info)
	})

	t.Run("successful info fetch", func(t *testing.T) {
		t.Parallel()
		user := &model.User{ID: 1, Username: "test", Balance: 100}
		inventory := []model.Inventory{{Type: "item", Quantity: 1}}
		history := &model.CoinHistory{
			Received: []model.CoinsReceived{{FromUser: "user2", Amount: 50}},
		}

		ts.repo.EXPECT().
			WithTx(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, fn func(repository.DB) error) error {
				return fn(nil)
			})
		ts.repo.EXPECT().FindUser(gomock.Any(), nil, user.Username).Return(user, nil)
		ts.repo.EXPECT().ListInventory(gomock.Any(), nil, user.ID).Return(inventory, nil)
		ts.repo.EXPECT().ListTransactions(gomock.Any(), nil, user.ID).Return(history, nil)

		info, err := ts.users.Info(ctx, user.Username)
		assert.NoError(t, err)
		assert.Equal(t, user.Balance, info.Coins)
		assert.Equal(t, inventory, info.Inventory)
		assert.Equal(t, history, info.CoinHistory)
	})
}

func TestService_Transfer(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("validation cases", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name          string
			from          string
			to            string
			amount        int
			expectedError error
		}{
			{
				name:          "empty sender",
				from:          "",
				to:            "user2",
				amount:        100,
				expectedError: model.ErrUnauthorized,
			},
			{
				name:          "empty receiver",
				from:          "user1",
				to:            "",
				amount:        100,
				expectedError: model.ErrBadRequest,
			},
			{
				name:          "zero amount",
				from:          "user1",
				to:            "user2",
				amount:        0,
				expectedError: model.ErrBadRequest,
			},
			{
				name:          "negative amount",
				from:          "user1",
				to:            "user2",
				amount:        -100,
				expectedError: model.ErrBadRequest,
			},
			{
				name:          "self transfer",
				from:          "user1",
				to:            "user1",
				amount:        100,
				expectedError: model.ErrBadRequest,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				ts := newTestSuite(t)
				err := ts.users.Transfer(ctx, tt.from, tt.to, tt.amount)
				assert.ErrorIs(t, err, tt.expectedError)
			})
		}
	})

	t.Run("successful transfer", func(t *testing.T) {
		t.Parallel()
		ts := newTestSuite(t)

		sender := &model.User{ID: 1, Username: "sender", Balance: 1000}
		receiver := &model.User{ID: 2, Username: "receiver", Balance: 500}
		amount := 100

		ts.repo.EXPECT().
			WithTx(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, fn func(repository.DB) error) error {
				return fn(nil)
			})

		ts.repo.EXPECT().
			FindUser(gomock.Any(), nil, sender.Username).
			Return(sender, nil)

		ts.repo.EXPECT().
			FindUser(gomock.Any(), nil, receiver.Username).
			Return(receiver, nil)

		ts.repo.EXPECT().
			MakeTransfer(gomock.Any(), nil, sender.ID, receiver.ID, amount).
			Return(nil)

		err := ts.users.Transfer(ctx, sender.Username, receiver.Username, amount)
		assert.NoError(t, err)
	})

	t.Run("insufficient funds", func(t *testing.T) {
		t.Parallel()
		ts := newTestSuite(t)

		sender := &model.User{ID: 1, Username: "sender", Balance: 50}
		amount := 100

		ts.repo.EXPECT().
			WithTx(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, fn func(repository.DB) error) error {
				return fn(nil)
			})

		ts.repo.EXPECT().
			FindUser(gomock.Any(), nil, sender.Username).
			Return(sender, nil)

		err := ts.users.Transfer(ctx, sender.Username, "receiver", amount)
		assert.ErrorIs(t, err, model.ErrInsufficientFunds)
	})
}
