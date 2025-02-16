package auth

import (
	"context"
	"database/sql"
	"errors"

	"github.com/esklo/avito-backend-winter-2025/internal/model"
	"github.com/esklo/avito-backend-winter-2025/internal/repository"
	"github.com/esklo/avito-backend-winter-2025/internal/service"
	"github.com/golang-jwt/jwt/v5"
)

var _ service.Authenticator = (*Service)(nil)

var (
	ErrInvalidToken         = errors.New("invalid token")
	ErrInvalidSigningMethod = errors.New("unexpected signing method")
)

type Service struct {
	repo   repository.Repository
	users  service.UserManager
	hasher service.Hasher
	secret []byte
}

func NewService(repo repository.Repository, users service.UserManager, hasher service.Hasher, secret []byte) *Service {
	return &Service{
		repo:   repo,
		users:  users,
		hasher: hasher,
		secret: secret,
	}
}

func (s *Service) Login(ctx context.Context, username, password string) (string, error) {
	if username == "" || password == "" {
		return "", model.ErrBadRequest
	}

	user, err := s.repo.FindUser(ctx, nil, username)
	if errors.Is(err, sql.ErrNoRows) {
		user, err = s.users.Create(ctx, username, password)
	}

	if err != nil {
		return "", model.ErrUnauthorized
	}

	if !s.hasher.Verify(password, user.Password, user.Salt) {
		return "", model.ErrUnauthorized
	}

	return createToken(newClaims(user.Username), s.secret)
}

func (s *Service) ValidateToken(_ context.Context, token string) (string, error) {
	if token == "" {
		return "", ErrInvalidToken
	}

	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidSigningMethod
		}

		return s.secret, nil
	})
	if err != nil {
		return "", ErrInvalidToken
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok || !jwtToken.Valid {
		return "", ErrInvalidToken
	}

	return claims.GetSubject()
}
