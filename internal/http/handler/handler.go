package handler

import "github.com/esklo/avito-backend-winter-2025/internal/service"

type CtxKey string

const CtxUsernameKey CtxKey = "username"

type Container interface {
	Users() service.UserManager
	Auth() service.Authenticator
	Shop() service.Shop
}
type Handler struct {
	container Container
}

func New(container Container) *Handler {
	return &Handler{container: container}
}
