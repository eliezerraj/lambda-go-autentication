package service

import (
	"context"
	"github.com/rs/zerolog/log"
	"time"

	"github.com/lambda-go-autentication/internal/core/domain"
	"github.com/lambda-go-autentication/internal/erro"
	"github.com/lambda-go-autentication/internal/repository"

	"github.com/golang-jwt/jwt/v4"
	"github.com/aws/aws-xray-sdk-go/xray"
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

func (a AuthService) Login(ctx context.Context, credential domain.Credential) (*domain.Authentication, error){
	childLogger.Debug().Msg("Login")

	_, root := xray.BeginSubsegment(ctx, "Service.Login")
	defer root.Close(nil)

	_, err := a.authRepository.Login(ctx, credential)
	if err != nil {
		return nil, err
	}

	// get scopes associated with a credential
	credential_scope, err := a.authRepository.QueryCredentialScope(ctx, credential)
	if err != nil {
		return nil, err
	}

	// Set a JWT expiration date 
	expirationTime := time.Now().Add(120 * time.Minute)

	// Create a JWT Oauth 2.0 with all scopes and expiration date
	jwtData := &domain.JwtData{
								Username: credential.User,
								Scope: credential_scope.Scope,
								RegisteredClaims: jwt.RegisteredClaims{
									ExpiresAt: jwt.NewNumericDate(expirationTime), 	// JWT expiry time is unix milliseconds
								},
	}

	// Add the claims and sign the token
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

	// Create a new credential
	_, err := a.authRepository.SignIn(credential)
	if err != nil {
		return nil, err
	}
	return &credential,nil
}

func (a AuthService) TokenValidation(credential domain.Credential) (bool, error){
	childLogger.Debug().Msg("TokenValidation")

	// Check with token is signed 
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

	// Save the credentials scopes
	_, err := a.authRepository.AddScope(credential_scope)
	if err != nil {
		return nil, err
	}

	return &credential_scope, nil
}

func (a AuthService) QueryCredentialScope(ctx context.Context, credential domain.Credential) (*domain.CredentialScope, error){
	childLogger.Debug().Msg("QueryCredentialScope")

	_, root := xray.BeginSubsegment(ctx, "Service.QueryCredentialScope")
	defer root.Close(nil)

	// Query all scope linked with the credentials
	credential_scope, err := a.authRepository.QueryCredentialScope(ctx, credential)
	if err != nil {
		return nil, err
	}

	return credential_scope, nil
}

func (a AuthService) RefreshToken(credential domain.Credential) (*domain.Authentication, error){
	childLogger.Debug().Msg("RefreshToken")

	// Check with token is signed 
	claims := &domain.JwtData{}
	tkn, err := jwt.ParseWithClaims(credential.Token, claims, func(token *jwt.Token) (interface{}, error) {
		return a.jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, erro.ErrStatusUnauthorized
		}
		return nil, erro.ErrTokenExpired
	}

	if !tkn.Valid {
		return nil, erro.ErrStatusUnauthorized
	}

	// Check if the token is still valid
	if time.Until(claims.ExpiresAt.Time) > (50 * time.Minute) {
		return nil, erro.ErrTokenStillValid
	}

	// Set a new tokens claims
	expirationTime := time.Now().Add(60 * time.Minute)
	claims.ExpiresAt = jwt.NewNumericDate(expirationTime)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(a.jwtKey)

	auth := domain.Authentication{	Token: tokenString, 
									ExpirationTime :expirationTime}

	return &auth,nil
}