package handler

import (
	"context"
	"net/http"

	"github.com/esklo/avito-backend-winter-2025/internal/http/render"
	"github.com/esklo/avito-backend-winter-2025/internal/model"
)

func (h *Handler) Info(w http.ResponseWriter, r *http.Request) {
	username, err := usernameFromCtx(r.Context())
	if err != nil {
		render.Error(w, err)

		return
	}

	info, err := h.container.Users().Info(r.Context(), username)
	if err != nil {
		render.Error(w, err)

		return
	}

	render.Success(w, info)
}

func usernameFromCtx(ctx context.Context) (string, error) {
	username, ok := ctx.Value(CtxUsernameKey).(string)
	if !ok || username == "" {
		return "", model.ErrUnauthorized
	}

	return username, nil
}
