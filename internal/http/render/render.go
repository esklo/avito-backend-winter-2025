package render

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/esklo/avito-backend-winter-2025/internal/model"
)

func Success(w http.ResponseWriter, data any) {
	render(w, http.StatusOK, data)
}

func Error(w http.ResponseWriter, err error) {
	code := getStatusCode(err)

	render(w, code, map[string]any{
		"errors": err.Error(),
	})
}

func render(w http.ResponseWriter, code int, resp any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		Error(w, err)
	}
}

func getStatusCode(err error) int {
	switch {
	case errors.Is(err, model.ErrBadRequest):
		return http.StatusBadRequest
	case errors.Is(err, model.ErrUnauthorized):
		return http.StatusUnauthorized
	case errors.Is(err, model.ErrNotFound):
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
