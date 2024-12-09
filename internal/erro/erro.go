package erro

import (
	"errors"
)

var (
	ErrCertRevoked = errors.New("unauthorized cert revoked")
	ErrParseCert = errors.New("unable to parse x509 cert")
	ErrDecodeCert = errors.New("failed to decode pem-encoded cert")
	ErrDecodeKey = errors.New("error decode rsa key")
	ErrTokenExpired	= errors.New("token expired")
	ErrStatusUnauthorized = errors.New("invalid Token")
	ErrArnMalFormad = errors.New("unauthorized arn scoped malformed")
	ErrBearTokenFormad = errors.New("unauthorized token not informed")
	ErrUnmarshal = errors.New("erro unmarshall")
	ErrInsert 	= errors.New("insert error")
	ErrPreparedQuery = errors.New("prepare dynamo query erro")
	ErrQuery = errors.New("query table error")
	ErrNotFound = errors.New("data not found")
	ErrList	= errors.New("list query error")
	ErrMethodNotAllowed	= errors.New("method not allowed")
	ErrQueryEmpty	= errors.New("query parameters missing")
	ErrTokenStillValid = errors.New("token is still valid")
)