package http

import (
	"context"
	"net/http"
	"strings"

	"github.com/esklo/avito-backend-winter-2025/internal/http/handler"
	"github.com/esklo/avito-backend-winter-2025/internal/http/render"
	"github.com/esklo/avito-backend-winter-2025/internal/model"
)

func (s *Server) withMiddlewares(next http.HandlerFunc) http.HandlerFunc {
	return s.withRecover(
		s.withCORS(next),
	)
}

func (s *Server) withAuth(next http.HandlerFunc) http.HandlerFunc {
	return s.withMiddlewares(func(w http.ResponseWriter, r *http.Request) {
		token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

		username, err := s.container.Auth().ValidateToken(r.Context(), token)
		if err != nil {
			render.Error(w, err)

			return
		}

		ctx := context.WithValue(r.Context(), handler.CtxUsernameKey, username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *Server) withRecover(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				render.Error(w, model.ErrInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	}
}

func (s *Server) withCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)

			return
		}

		next.ServeHTTP(w, r)
	}
}
