package service

import (
	"github.com/rs/zerolog/log"
	"time"

	"github.com/lambda-go-autentication/internal/core/domain"
	"github.com/lambda-go-autentication/internal/erro"

	"github.com/golang-jwt/jwt/v4"
)

var childLogger = log.With().Str("service", "AuthService").Logger()

type AuthService struct {
	jwtKey	[]byte
}

func NewAuthService(	//cardRepository repository.CardRepository,
						jwtKey []byte  ) *AuthService{
	childLogger.Debug().Msg("NewAuthService")
	return &AuthService{
		jwtKey: jwtKey,
	}
}

func (a AuthService) Login(credential domain.Credential) (*domain.Authentication, error){
	childLogger.Debug().Msg("Login")

	expirationTime := time.Now().Add(5 * time.Minute)

	jwtData := &domain.JwtData{
								Username: credential.User,
								Scope: "query-product",
								RegisteredClaims: jwt.RegisteredClaims{
								ExpiresAt: jwt.NewNumericDate(expirationTime), 	// JWT expiry time is unix milliseconds
								},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtData)
	tokenString, err := token.SignedString(a.jwtKey)
	if err != nil {
		return nil, err
	}

	auth := domain.Authentication{	Token: tokenString, 
									ExpirationTime :expirationTime}

	return &auth,nil
}

func (a AuthService) TokenValidation(credential domain.Credential) (bool, error){
	childLogger.Debug().Msg("TokenValidation")

	claims := &domain.JwtData{}
	tkn, err := jwt.ParseWithClaims(credential.Token, claims, func(token *jwt.Token) (interface{}, error) {
		return a.jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return false, erro.ErrStatusUnauthorized
		}
		return false, erro.ErrTokenExpired
	}

	if !tkn.Valid {
		return false, erro.ErrStatusUnauthorized
	}

	return true ,nil
}