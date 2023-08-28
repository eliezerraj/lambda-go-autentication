package erro

import (
	"errors"

)

var (
	ErrStatusUnauthorized 	= errors.New("Invalid Token")
	ErrTokenExpired		 	= errors.New("Token expired")
	ErrBadRequest		 	= errors.New("Internal error")
	ErrUnmarshal			= errors.New("Erro Unmarshall")
	ErrMethodNotAllowed		= errors.New("Method not allowed")
)