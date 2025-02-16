package di

import (
	"testing"

	"github.com/esklo/avito-backend-winter-2025/internal/config"
	"github.com/esklo/avito-backend-winter-2025/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("successful container creation", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		mockRepo := mocks.NewMockRepository(ctrl)
		cfg := &config.Config{
			App: config.AppConfig{
				JWTSecret: []byte("test-secret"),
			},
		}

		container := New(cfg, mockRepo)
		require.NotNil(t, container)

		assert.Equal(t, cfg, container.Config())
		assert.NotNil(t, container.Log())
		assert.NotNil(t, container.Hasher())
		assert.NotNil(t, container.Auth())
		assert.NotNil(t, container.Users())
		assert.NotNil(t, container.Shop())
	})
}

func TestContainer_Dependencies(t *testing.T) {
	t.Parallel()

	t.Run("verify dependency injection", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		mockRepo := mocks.NewMockRepository(ctrl)
		cfg := &config.Config{
			App: config.AppConfig{
				JWTSecret: []byte("test-secret"),
			},
		}

		container := New(cfg, mockRepo)

		assert.NotNil(t, container.Config())
		assert.Equal(t, cfg, container.Config())

		assert.NotNil(t, container.Log())

		hasher := container.Hasher()
		assert.NotNil(t, hasher)

		users := container.Users()
		assert.NotNil(t, users)

		auth := container.Auth()
		assert.NotNil(t, auth)

		shop := container.Shop()
		assert.NotNil(t, shop)
	})
}
