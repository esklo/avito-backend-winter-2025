package handler

import (
	"encoding/json"
	"net/http"

	"github.com/esklo/avito-backend-winter-2025/internal/http/render"
	"github.com/esklo/avito-backend-winter-2025/internal/model"
)

type transferRequest struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}

func (h *Handler) Transfer(w http.ResponseWriter, r *http.Request) {
	username, err := usernameFromCtx(r.Context())
	if err != nil {
		render.Error(w, err)

		return
	}

	var req transferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Error(w, model.ErrBadRequest)

		return
	}

	err = h.container.Users().Transfer(r.Context(), username, req.ToUser, req.Amount)
	if err != nil {
		render.Error(w, err)

		return
	}

	render.Success(w, nil)
}
