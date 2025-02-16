package model

import "errors"

var (
	ErrBadRequest          = errors.New("bad request")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrInternalServerError = errors.New("internal server error")
	ErrInsufficientFunds   = errors.New("insufficient funds")
	ErrNotFound            = errors.New("not found")
)
