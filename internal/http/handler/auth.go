package handler

import (
	"encoding/json"
	"net/http"

	"github.com/esklo/avito-backend-winter-2025/internal/http/render"
	"github.com/esklo/avito-backend-winter-2025/internal/model"
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Error(w, model.ErrBadRequest)

		return
	}

	token, err := h.container.Auth().Login(r.Context(), req.Username, req.Password)
	if err != nil {
		render.Error(w, err)

		return
	}

	render.Success(w, loginResponse{Token: token})
}
