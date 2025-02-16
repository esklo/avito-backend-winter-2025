package render

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/esklo/avito-backend-winter-2025/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSuccess(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		data interface{}
	}{
		{
			name: "nil data",
			data: nil,
		},
		{
			name: "string data",
			data: "test",
		},
		{
			name: "map data",
			data: map[string]string{"key": "value"},
		},
		{
			name: "struct data",
			data: struct {
				Field string `json:"field"`
			}{Field: "value"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			w := httptest.NewRecorder()

			Success(w, tt.data)

			assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

			if tt.data != nil {
				var response interface{}
				err := json.NewDecoder(w.Body).Decode(&response)
				require.NoError(t, err)
			}
		})
	}
}

func TestError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		err           error
		expectedCode  int
		expectedError string
	}{
		{
			name:          "bad request error",
			err:           model.ErrBadRequest,
			expectedCode:  http.StatusBadRequest,
			expectedError: model.ErrBadRequest.Error(),
		},
		{
			name:          "unauthorized error",
			err:           model.ErrUnauthorized,
			expectedCode:  http.StatusUnauthorized,
			expectedError: model.ErrUnauthorized.Error(),
		},
		{
			name:          "not found error",
			err:           model.ErrNotFound,
			expectedCode:  http.StatusNotFound,
			expectedError: model.ErrNotFound.Error(),
		},
		{
			name:          "internal server error",
			err:           errors.New("unexpected error"),
			expectedCode:  http.StatusInternalServerError,
			expectedError: "unexpected error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			w := httptest.NewRecorder()

			Error(w, tt.err)

			assert.Equal(t, tt.expectedCode, w.Code)
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

			var response struct {
				Errors string `json:"errors"`
			}
			err := json.NewDecoder(w.Body).Decode(&response)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedError, response.Errors)
		})
	}
}

func TestGetStatusCode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		err          error
		expectedCode int
	}{
		{
			name:         "bad request error",
			err:          model.ErrBadRequest,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "unauthorized error",
			err:          model.ErrUnauthorized,
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "not found error",
			err:          model.ErrNotFound,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "wrapped error",
			err:          errors.Join(model.ErrBadRequest, errors.New("context")),
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "unknown error",
			err:          errors.New("unknown"),
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			code := getStatusCode(tt.err)
			assert.Equal(t, tt.expectedCode, code)
		})
	}
}
