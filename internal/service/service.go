package service

import (
	"context"
	"github.com/rs/zerolog/log"
	"time"

	"github.com/lambda-go-autentication/internal/config/observability"
	"github.com/lambda-go-autentication/internal/core"
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

func (a AuthService) Login(ctx context.Context, credential core.Credential) (*core.Authentication, error){
	childLogger.Debug().Msg("Login")
	childLogger.Debug().Interface("credential :",credential).Msg("")

	span := observability.Span(ctx, "service.login")	
    defer span.End()

	_, err := a.authRepository.Login(ctx, credential)
	if err != nil {
		return nil, err
	}

	// get scopes associated with a credential
	credential_scope, err := a.authRepository.QueryCredentialScope(ctx, credential)
	if err != nil {
		return nil, err
	}

	span_jwt := observability.Span(ctx, "service.create_jwt")	
    
	// Set a JWT expiration date 
	expirationTime := time.Now().Add(720 * time.Minute)

	// Create a JWT Oauth 2.0 with all scopes and expiration date
	jwtData := &core.JwtData{
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
	defer span_jwt.End()
	
	auth := core.Authentication{	Token: tokenString, 
									ExpirationTime :expirationTime}

	return &auth,nil
}

func (a AuthService) SignIn(ctx context.Context, credential core.Credential) (*core.Credential, error){
	childLogger.Debug().Msg("SignIn")

	span := observability.Span(ctx, "service.signIn")	
    defer span.End()

	// Create a new credential
	res, err := a.authRepository.SignIn(ctx, credential)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (a AuthService) TokenValidation(ctx context.Context, credential core.Credential) (bool, error){
	childLogger.Debug().Msg("TokenValidation")

	span := observability.Span(ctx, "service.tokenValidation")	
    defer span.End()

	// Check with token is signed 
	claims := &core.JwtData{}
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

func (a AuthService) AddScope(ctx context.Context, credential_scope core.CredentialScope) (*core.CredentialScope, error){
	childLogger.Debug().Msg("AddScope")

	span := observability.Span(ctx, "service.addScope")	
    defer span.End()

	// Save the credentials scopes
	res, err := a.authRepository.AddScope(ctx, credential_scope)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (a AuthService) QueryCredentialScope(ctx context.Context, credential core.Credential) (*core.CredentialScope, error){
	childLogger.Debug().Msg("QueryCredentialScope")

	span := observability.Span(ctx, "service.queryCredentialScope")	
    defer span.End()

	// Query all scope linked with the credentials
	credential_scope, err := a.authRepository.QueryCredentialScope(ctx, credential)
	if err != nil {
		return nil, err
	}

	return credential_scope, nil
}

func (a AuthService) RefreshToken(ctx context.Context, credential core.Credential) (*core.Authentication, error){
	childLogger.Debug().Msg("RefreshToken")

	span := observability.Span(ctx, "service.refreshToken")	
    defer span.End()

	// Check with token is signed 
	claims := &core.JwtData{}
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

	auth := core.Authentication{	Token: tokenString, 
									ExpirationTime :expirationTime}

	return &auth,nil
}