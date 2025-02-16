package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/esklo/avito-backend-winter-2025/internal/model"

	"github.com/esklo/avito-backend-winter-2025/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type testSuite struct {
	handler   *Handler
	container *mocks.MockContainer
	shop      *mocks.MockShop
	users     *mocks.MockUserManager
	auth      *mocks.MockAuthenticator
}

func newTestSuite(t *testing.T) *testSuite {
	ctrl := gomock.NewController(t)

	container := mocks.NewMockContainer(ctrl)
	shop := mocks.NewMockShop(ctrl)
	users := mocks.NewMockUserManager(ctrl)
	auth := mocks.NewMockAuthenticator(ctrl)

	container.EXPECT().Shop().Return(shop).AnyTimes()
	container.EXPECT().Users().Return(users).AnyTimes()
	container.EXPECT().Auth().Return(auth).AnyTimes()

	return &testSuite{
		container: container,
		shop:      shop,
		users:     users,
		auth:      auth,
		handler:   New(container),
	}
}

func TestHandler_Login(t *testing.T) {
	t.Parallel()

	t.Run("successful login", func(t *testing.T) {
		t.Parallel()
		ts := newTestSuite(t)

		req := loginRequest{
			Username: "user",
			Password: "pass",
		}

		ts.auth.EXPECT().
			Login(gomock.Any(), req.Username, req.Password).
			Return("test-token", nil)

		w := httptest.NewRecorder()
		data, err := json.Marshal(req)
		require.NoError(t, err)

		r := httptest.NewRequest(http.MethodPost, "/api/auth", bytes.NewReader(data))
		ts.handler.Login(w, r)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp loginResponse
		err = json.NewDecoder(w.Body).Decode(&resp)
		require.NoError(t, err)
		assert.Equal(t, "test-token", resp.Token)
	})
}

func TestHandler_Buy(t *testing.T) {
	t.Parallel()

	t.Run("successful purchase", func(t *testing.T) {
		t.Parallel()
		ts := newTestSuite(t)

		ts.shop.EXPECT().
			BuyItem(gomock.Any(), "hoody", "test-user").
			Return(nil)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/api/buy/hoody", nil)
		r.SetPathValue("name", "hoody")
		ctx := context.WithValue(r.Context(), CtxUsernameKey, "test-user")
		r = r.WithContext(ctx)

		ts.handler.Buy(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("unauthorized", func(t *testing.T) {
		t.Parallel()
		ts := newTestSuite(t)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/api/buy/hoody", nil)

		ts.handler.Buy(w, r)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestHandler_Transfer(t *testing.T) {
	t.Parallel()

	t.Run("successful transfer", func(t *testing.T) {
		t.Parallel()
		ts := newTestSuite(t)

		req := transferRequest{
			ToUser: "receiver",
			Amount: 100,
		}

		ts.users.EXPECT().
			Transfer(gomock.Any(), "sender", req.ToUser, req.Amount).
			Return(nil)

		w := httptest.NewRecorder()
		data, err := json.Marshal(req)
		require.NoError(t, err)

		r := httptest.NewRequest(http.MethodPost, "/api/sendCoin", bytes.NewReader(data))
		ctx := context.WithValue(r.Context(), CtxUsernameKey, "sender")
		r = r.WithContext(ctx)

		ts.handler.Transfer(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("unauthorized", func(t *testing.T) {
		t.Parallel()
		ts := newTestSuite(t)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/sendCoin", nil)

		ts.handler.Transfer(w, r)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestHandler_Info(t *testing.T) {
	t.Parallel()

	t.Run("get user info", func(t *testing.T) {
		t.Parallel()
		ts := newTestSuite(t)

		ts.users.EXPECT().
			Info(gomock.Any(), "sender").
			Return(&model.Info{}, nil)

		w := httptest.NewRecorder()

		r := httptest.NewRequest(http.MethodGet, "/api/info", nil)
		ctx := context.WithValue(r.Context(), CtxUsernameKey, "sender")
		r = r.WithContext(ctx)

		ts.handler.Info(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("unauthorized", func(t *testing.T) {
		t.Parallel()
		ts := newTestSuite(t)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/api/info", nil)

		ts.handler.Info(w, r)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
