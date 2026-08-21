package entity

import "errors"

var (
	ErrInvalidParams = errors.New("invalid params")
	ErrNotFound      = errors.New("not found")
)
