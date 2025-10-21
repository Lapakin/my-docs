package jwt

import (
	"errors"
)

var (
	ErrBadJwt        = errors.New("JWT token invalid")
	ErrUnmarshal     = errors.New("cannot unmarshal")
	ErrJwtEncode     = errors.New("error while JWT encoding")
	ErrRoleMissing   = errors.New("role is missing")
	ErrUserIDMissing = errors.New("user id is missing")
	ErrNameMissing   = errors.New("user name is missing")
)
