package handler

import (
	"net/http"

	"github.com/esklo/avito-backend-winter-2025/internal/http/render"
)

func (h *Handler) Buy(w http.ResponseWriter, r *http.Request) {
	username, err := usernameFromCtx(r.Context())
	if err != nil {
		render.Error(w, err)

		return
	}

	err = h.container.Shop().BuyItem(r.Context(), r.PathValue("name"), username)
	if err != nil {
		render.Error(w, err)

		return
	}

	render.Success(w, nil)
}
