package service

import (
	"github.com/rs/zerolog/log"
	"time"

	"github.com/lambda-go-autentication/internal/core/domain"
	"github.com/lambda-go-autentication/internal/erro"
	"github.com/lambda-go-autentication/internal/repository"

	"github.com/golang-jwt/jwt/v4"
)

var childLogger = log.With().Str("service", "AuthService").Logger()

type AuthService struct {
	jwtKey			[]byte
	authRepository	*repository.AuthRepository
}

func NewAuthService(jwtKey []byte,
					authRepository *repository.AuthRepository) *AuthService{
	childLogger.Debug().Msg("NewAuthService")
	return &AuthService{
		jwtKey: jwtKey,
		authRepository: authRepository,
	}
}

func (a AuthService) Login(credential domain.Credential) (*domain.Authentication, error){
	childLogger.Debug().Msg("Login")

	_, err := a.authRepository.Login(credential)
	if err != nil {
		return nil, err
	}
	//childLogger.Debug().Interface(credential_login).Msg("Login")

	expirationTime := time.Now().Add(60 * time.Minute)

	user_scope := []string{"info.read",
							"a.read",
							"sum.write",
							"version",
							"header.read"}

	jwtData := &domain.JwtData{
								Username: credential.User,
								Scope: user_scope,
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

func (a AuthService) SignIn(credential domain.Credential) (*domain.Credential, error){
	childLogger.Debug().Msg("SignIn")

	_, err := a.authRepository.SignIn(credential)
	if err != nil {
		return nil, err
	}
	return &credential,nil
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

func (a AuthService) AddScope(credential_scope domain.CredentialScope) (*domain.CredentialScope, error){
	childLogger.Debug().Msg("AddScope")

	_, err := a.authRepository.AddScope(credential_scope)
	if err != nil {
		return nil, err
	}

	return &credential_scope, nil
}

func (a AuthService) QueryCredentialScope(credential domain.Credential) (*domain.CredentialScope, error){
	childLogger.Debug().Msg("QueryCredentialScope")

	credential_scope, err := a.authRepository.QueryCredentialScope(credential)
	if err != nil {
		return nil, err
	}

	return credential_scope, nil
}