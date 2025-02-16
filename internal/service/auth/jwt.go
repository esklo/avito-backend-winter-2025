package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// todo: should be moved to config
const (
	issuer   = "shop"
	tokenExp = 24 * time.Hour
)

func newClaims(username string) jwt.RegisteredClaims {
	now := time.Now()

	return jwt.RegisteredClaims{
		Issuer:    issuer,
		Subject:   username,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(tokenExp)),
	}
}

func createToken(claims jwt.Claims, secret []byte) (string, error) {
	return jwt.
		NewWithClaims(
			jwt.SigningMethodHS256,
			claims,
		).
		SignedString(secret)
}
