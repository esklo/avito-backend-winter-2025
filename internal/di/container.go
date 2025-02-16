package di

import (
	"log/slog"

	"github.com/esklo/avito-backend-winter-2025/internal/service"

	"github.com/esklo/avito-backend-winter-2025/internal/config"
	"github.com/esklo/avito-backend-winter-2025/internal/repository"
	"github.com/esklo/avito-backend-winter-2025/internal/service/auth"
	"github.com/esklo/avito-backend-winter-2025/internal/service/hasher"
	"github.com/esklo/avito-backend-winter-2025/internal/service/shop"
	"github.com/esklo/avito-backend-winter-2025/internal/service/user"
)

type Container struct {
	cfg  *config.Config
	repo repository.Repository

	log *slog.Logger

	hasher *hasher.Argon2
	auth   *auth.Service
	users  *user.Service
	shop   *shop.Service
}

func New(cfg *config.Config, repo repository.Repository) *Container {
	c := &Container{
		cfg:    cfg,
		repo:   repo,
		log:    slog.Default(),
		hasher: hasher.NewArgon2(),
	}
	c.initServices()

	return c
}

func (c *Container) initServices() {
	c.users = user.NewService(c.repo, c.hasher)
	c.auth = auth.NewService(c.repo, c.users, c.hasher, c.cfg.App.JWTSecret)
	c.shop = shop.NewService(c.repo)
}

func (c *Container) Config() *config.Config      { return c.cfg }
func (c *Container) Log() *slog.Logger           { return c.log }
func (c *Container) Hasher() service.Hasher      { return c.hasher }
func (c *Container) Auth() service.Authenticator { return c.auth }
func (c *Container) Users() service.UserManager  { return c.users }
func (c *Container) Shop() service.Shop          { return c.shop }
