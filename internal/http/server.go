package http

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/esklo/avito-backend-winter-2025/internal/config"

	"github.com/esklo/avito-backend-winter-2025/internal/service"

	"github.com/esklo/avito-backend-winter-2025/internal/http/handler"
)

//go:generate mockgen -destination=../../mocks/mock_server_container.go -package=mocks github.com/esklo/avito-backend-winter-2025/internal/http Container
type Container interface {
	Log() *slog.Logger
	Config() *config.Config
	Users() service.UserManager
	Auth() service.Authenticator
	Shop() service.Shop
}

type Server struct {
	srv       *http.Server
	router    *http.ServeMux
	container Container
}

func NewServer(container Container) *Server {
	s := &Server{
		router:    http.NewServeMux(),
		container: container,
	}

	s.setupRoutes()

	s.srv = &http.Server{
		Addr:              container.Config().HTTP.Address(),
		Handler:           s.router,
		ReadHeaderTimeout: 3 * time.Second,
	}

	return s
}

func (s *Server) Run(ctx context.Context) error {
	s.container.Log().Info("listening on", "addr", s.srv.Addr)

	defer func() { _ = s.Shutdown(ctx) }()

	return s.srv.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

func (s *Server) setupRoutes() {
	h := handler.New(s.container)

	s.router.Handle("POST /api/auth", s.withMiddlewares(h.Login))
	s.router.Handle("GET /api/info", s.withAuth(h.Info))
	s.router.Handle("GET /api/buy/{name}", s.withAuth(h.Buy))
	s.router.Handle("POST /api/sendCoin", s.withAuth(h.Transfer))
}
