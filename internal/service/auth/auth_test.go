package auth

import (
	"context"
	"database/sql"
	"testing"

	"github.com/esklo/avito-backend-winter-2025/internal/model"
	"github.com/esklo/avito-backend-winter-2025/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type testSuite struct {
	auth   *Service
	repo   *mocks.MockRepository
	users  *mocks.MockUserManager
	hasher *mocks.MockHasher
}

func newTestSuite(t *testing.T) *testSuite {
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockRepository(ctrl)
	users := mocks.NewMockUserManager(ctrl)
	hasher := mocks.NewMockHasher(ctrl)

	return &testSuite{
		repo:   repo,
		users:  users,
		hasher: hasher,
		auth:   NewService(repo, users, hasher, []byte("test-secret")),
	}
}

func TestService_Login(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("validation cases", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name          string
			username      string
			password      string
			expectedError error
		}{
			{
				name:          "empty both",
				expectedError: model.ErrBadRequest,
			},
			{
				name:          "empty password",
				username:      "user",
				expectedError: model.ErrBadRequest,
			},
			{
				name:          "empty username",
				password:      "pass",
				expectedError: model.ErrBadRequest,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()
				ts := newTestSuite(t)

				token, err := ts.auth.Login(ctx, tt.username, tt.password)
				assert.ErrorIs(t, err, tt.expectedError)
				assert.Empty(t, token)
			})
		}
	})

	t.Run("existing user login", func(t *testing.T) {
		t.Parallel()
		ts := newTestSuite(t)

		user := &model.User{
			ID:       1,
			Username: "user",
			Password: []byte("hashed"),
			Salt:     []byte("salt"),
		}

		ts.repo.EXPECT().
			FindUser(gomock.Any(), nil, user.Username).
			Return(user, nil)

		ts.hasher.EXPECT().
			Verify("password", user.Password, user.Salt).
			Return(true)

		token, err := ts.auth.Login(ctx, user.Username, "password")
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("new user registration", func(t *testing.T) {
		t.Parallel()
		ts := newTestSuite(t)

		user := &model.User{
			ID:       1,
			Username: "user",
			Password: []byte("hashed"),
			Salt:     []byte("salt"),
		}

		ts.repo.EXPECT().
			FindUser(gomock.Any(), nil, user.Username).
			Return(nil, sql.ErrNoRows)

		ts.users.EXPECT().
			Create(gomock.Any(), user.Username, "password").
			Return(user, nil)

		ts.hasher.EXPECT().
			Verify("password", user.Password, user.Salt).
			Return(true)

		token, err := ts.auth.Login(ctx, user.Username, "password")
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})
}

func TestService_ValidateToken(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("empty token", func(t *testing.T) {
		t.Parallel()
		ts := newTestSuite(t)

		username, err := ts.auth.ValidateToken(ctx, "")
		assert.ErrorIs(t, err, ErrInvalidToken)
		assert.Empty(t, username)
	})

	t.Run("invalid token", func(t *testing.T) {
		t.Parallel()
		ts := newTestSuite(t)

		username, err := ts.auth.ValidateToken(ctx, "invalid.token.here")
		assert.ErrorIs(t, err, ErrInvalidToken)
		assert.Empty(t, username)
	})

	t.Run("valid token", func(t *testing.T) {
		t.Parallel()
		ts := newTestSuite(t)

		token, err := createToken(newClaims("user"), []byte("test-secret"))
		require.NoError(t, err)

		username, err := ts.auth.ValidateToken(ctx, token)
		assert.NoError(t, err)
		assert.Equal(t, "user", username)
	})
}
