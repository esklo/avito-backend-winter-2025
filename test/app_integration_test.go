//go:build integration

package test

import (
	"context"
	"testing"
	"time"

	"github.com/esklo/avito-backend-winter-2025/internal/model"
	"github.com/esklo/avito-backend-winter-2025/internal/repository"
	"github.com/esklo/avito-backend-winter-2025/internal/service"
	"github.com/esklo/avito-backend-winter-2025/internal/service/auth"
	"github.com/esklo/avito-backend-winter-2025/internal/service/hasher"
	"github.com/esklo/avito-backend-winter-2025/internal/service/shop"
	"github.com/esklo/avito-backend-winter-2025/internal/service/user"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type testSuite struct {
	shop   service.Shop
	auth   service.Authenticator
	users  service.UserManager
	hasher service.Hasher
	repo   repository.Repository

	cleanup func()
}

func setupTestDB(t *testing.T) (*pgxpool.Pool, func()) {
	ctx := context.Background()
	container, err := postgres.Run(ctx,
		"postgres:13",
		postgres.WithInitScripts("../migrations/init.sql"),
		postgres.WithDatabase("test_db"),
		postgres.WithUsername("test_user"),
		postgres.WithPassword("test_password"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second),
		),
	)
	require.NoError(t, err)

	cleanup := func() {
		require.NoError(t, container.Terminate(ctx))
	}

	conn, err := container.ConnectionString(ctx)
	require.NoError(t, err)

	db, err := pgxpool.New(ctx, conn)
	require.NoError(t, err)

	return db, cleanup
}

func newTestSuite(t *testing.T) *testSuite {
	db, cleanup := setupTestDB(t)

	ts := &testSuite{
		hasher:  hasher.NewArgon2(),
		repo:    repository.New(db),
		cleanup: cleanup,
	}

	ts.shop = shop.NewService(ts.repo)
	ts.users = user.NewService(ts.repo, ts.hasher)
	ts.auth = auth.NewService(ts.repo, ts.users, ts.hasher, []byte("test-secret"))

	return ts
}

func (ts *testSuite) createTestUser(t *testing.T, username string) *model.User {
	_, err := ts.auth.Login(context.Background(), username, "test_password")
	require.NoError(t, err)

	u, err := ts.repo.FindUser(context.Background(), nil, username)
	require.NoError(t, err)

	return u
}

func TestAuthenticatorIntegration(t *testing.T) {
	t.Parallel()

	suite := newTestSuite(t)
	t.Cleanup(suite.cleanup)

	t.Run("authentication", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		_, err := suite.auth.Login(ctx, "user", "password")
		assert.NoError(t, err)

		_, err = suite.auth.Login(ctx, "user", "password1")
		assert.ErrorIs(t, err, model.ErrUnauthorized)

		_, err = suite.auth.Login(ctx, "user", "password")
		assert.NoError(t, err)
	})
}

func TestShopIntegration(t *testing.T) {
	t.Parallel()

	suite := newTestSuite(t)
	t.Cleanup(suite.cleanup)

	t.Run("purchase flow", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		u := suite.createTestUser(t, "buyer")

		infoBefore, err := suite.users.Info(ctx, u.Username)
		require.NoError(t, err)

		err = suite.shop.BuyItem(ctx, "pink-hoody", u.Username)
		require.NoError(t, err)

		infoAfter, err := suite.users.Info(ctx, u.Username)
		require.NoError(t, err)

		assert.Equal(t, infoBefore.Coins-500, infoAfter.Coins)
		assert.Contains(t, infoAfter.Inventory, model.Inventory{
			Type:     "pink-hoody",
			Quantity: 1,
		})
	})
}

func TestTransactionIntegration(t *testing.T) {
	t.Parallel()

	suite := newTestSuite(t)
	t.Cleanup(suite.cleanup)

	t.Run("transfer between users", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		user1 := suite.createTestUser(t, "sender")
		user2 := suite.createTestUser(t, "receiver")

		err := suite.users.Transfer(ctx, user1.Username, user2.Username, 100)
		require.NoError(t, err)

		info1, err := suite.users.Info(ctx, user1.Username)
		require.NoError(t, err)
		info2, err := suite.users.Info(ctx, user2.Username)
		require.NoError(t, err)

		assert.Equal(t, 900, info1.Coins)
		assert.Equal(t, 1100, info2.Coins)

		assert.Contains(t, info1.CoinHistory.Sent, model.CoinsSent{
			ToUser: user2.Username,
			Amount: 100,
		})
		assert.Contains(t, info2.CoinHistory.Received, model.CoinsReceived{
			FromUser: user1.Username,
			Amount:   100,
		})
	})
}
