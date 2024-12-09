package jwt

import (
	"fmt"
	"time"
	"context"
	"crypto/x509"
	"crypto/rsa"
    "encoding/pem"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"

	"github.com/rs/zerolog/log"

	"github.com/lambda-go-autentication/pkg/observability"
	"github.com/lambda-go-autentication/internal/erro"
	"github.com/lambda-go-autentication/internal/model"
)

var childLogger = log.With().Str("usecase", "jwt").Logger()

type UseCaseJwt struct{
	JwtKey		*string
	key_rsa_priv *rsa.PrivateKey
	key_rsa_pub *rsa.PublicKey
}

func NewUseCaseJwt(	jwtKey *string,
					key_rsa_priv *string,
					key_rsa_pub *string) *UseCaseJwt{
	childLogger.Debug().Msg("NewUseCaseJwt")

	_key_rsa_priv, err := ParsePemToRSAPriv(key_rsa_priv)
	if err != nil{
		childLogger.Error().Err(err).Msg("erro ParsePemToRSA !!!!")
	}
	_key_rsa_pub, err := ParsePemToRSAPub(key_rsa_pub)
	if err != nil{
		childLogger.Error().Err(err).Msg("erro ParsePemToRSA !!!!")
	}

	return &UseCaseJwt{
		JwtKey: jwtKey,
		key_rsa_priv: _key_rsa_priv,
		key_rsa_pub: _key_rsa_pub,
	}
}

func ParsePemToRSAPriv(private_key *string) (*rsa.PrivateKey, error){
	childLogger.Debug().Msg("ParsePemToRSA")

	block, _ := pem.Decode([]byte(*private_key))
	if block == nil || block.Type != "PRIVATE KEY" {
		childLogger.Error().Err(erro.ErrDecodeKey).Msg("erro Decode")
		return nil, erro.ErrDecodeKey
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		childLogger.Error().Err(err).Msg("erro ParsePKCS8PrivateKey")
		return nil, err
	}

	key_rsa := privateKey.(*rsa.PrivateKey)

	return key_rsa, nil
}

func ParsePemToRSAPub(public_key *string) (*rsa.PublicKey, error){
	childLogger.Debug().Msg("ParsePemToRSA")

	block, _ := pem.Decode([]byte(*public_key))
	if block == nil || block.Type != "PUBLIC KEY" {
		childLogger.Error().Err(erro.ErrDecodeKey).Msg("erro Decode")
		return nil, erro.ErrDecodeKey
	}

	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		childLogger.Error().Err(err).Msg("erro ParsePKCS8PrivateKey")
		return nil, err
	}

	key_rsa := pubInterface.(*rsa.PublicKey)

	return key_rsa, nil
}

func (u *UseCaseJwt) OAUTHToken(ctx context.Context, 
								credential model.Credential,
								credential_scope model.CredentialScope) (*model.Authentication, error){
	childLogger.Debug().Msg("OAUTHToken")

	childLogger.Debug().Interface("credential :",credential).Msg("")
	childLogger.Debug().Interface("credential_scope :",credential_scope).Msg("")
	childLogger.Debug().Interface("u.JwtKey :", u.JwtKey).Msg("")

	span := observability.Span(ctx, "usecase.OAUTHToken")
	defer span.End()

	// Set a JWT expiration date 
	expirationTime := time.Now().Add(720 * time.Minute)

	newUUID := uuid.New()
	uuidString := newUUID.String()

	// Create a JWT Oauth 2.0 with all scopes and expiration date
	jwtData := &model.JwtData{
								Username: credential.User,
								Scope: credential_scope.Scope,
								ISS: "lambda-go-autentication",
								Version: "2",
								JwtId: uuidString,
								TokenUse: "access",
								RegisteredClaims: jwt.RegisteredClaims{
									ExpiresAt: jwt.NewNumericDate(expirationTime), 	// JWT expiry time is unix milliseconds
								},
	}

	// Add the claims and sign the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtData)
	tokenString, err := token.SignedString([]byte(*u.JwtKey))
	if err != nil {
		return nil, err
	}
	
	auth := model.Authentication{Token: tokenString, 
								ExpirationTime :expirationTime}	

	return &auth ,nil
}

func (u *UseCaseJwt) TokenValidation(ctx context.Context, bearerToken string) (bool, error){
	childLogger.Debug().Msg("TokenValidation")

	span := observability.Span(ctx, "useCase.TokenValidation")	
    defer span.End()

	log.Debug().Interface("bearerToken : ", bearerToken).Msg("")

	claims := &model.JwtData{}
	tkn, err := jwt.ParseWithClaims(bearerToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(*u.JwtKey), nil
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

	return true, nil
}

func (u *UseCaseJwt) RefreshToken(ctx context.Context, bearerToken string) (*model.Authentication, error){
	childLogger.Debug().Msg("RefreshToken")

	span := observability.Span(ctx, "useCase.RefreshToken")	
    defer span.End()

	// Check with token is signed 
	claims := &model.JwtData{}
	tkn, err := jwt.ParseWithClaims(bearerToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(*u.JwtKey), nil
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
	if time.Until(claims.ExpiresAt.Time) > (719 * time.Minute) {
		return nil, erro.ErrTokenStillValid
	}

	// Set a new tokens claims
	expirationTime := time.Now().Add(720 * time.Minute)
	claims.ExpiresAt = jwt.NewNumericDate(expirationTime)
	claims.ISS = "lambda-go-autentication-refreshed"

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(*u.JwtKey))
	if err != nil {
		return nil, err
	}

	auth := model.Authentication{	Token: tokenString, 
									ExpirationTime :expirationTime}

	return &auth,nil
}
//----------------------------------------------------
func (u *UseCaseJwt) OAUTHTokenRSA(	ctx context.Context, 
									credential model.Credential,
									credential_scope model.CredentialScope) (*model.Authentication, error){
	childLogger.Debug().Msg("OAUTHTokenRSA")

	childLogger.Debug().Interface("credential :",credential).Msg("")
	childLogger.Debug().Interface("credential_scope :",credential_scope).Msg("")
	childLogger.Debug().Interface("u.key_rsa_priv :", u.key_rsa_priv).Msg("")
	childLogger.Debug().Interface("u.key_rsa_pub :", u.key_rsa_pub).Msg("")

	span := observability.Span(ctx, "usecase.OAUTHTokenRSA")
	defer span.End()

	// Set a JWT expiration date 
	expirationTime := time.Now().Add(720 * time.Minute)

	newUUID := uuid.New()
	uuidString := newUUID.String()

	// Create a JWT Oauth 2.0 with all scopes and expiration date
	jwtData := &model.JwtData{
								Username: credential.User,
								Scope: credential_scope.Scope,
								ISS: "lambda-go-autentication",
								Version: "2",
								JwtId: uuidString,
								TokenUse: "access-rsa",
								RegisteredClaims: jwt.RegisteredClaims{
									ExpiresAt: jwt.NewNumericDate(expirationTime), 	// JWT expiry time is unix milliseconds
								},
	}

	// Add the claims and sign the token
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwtData)
	tokenString, err := token.SignedString(u.key_rsa_priv)
	if err != nil {
		return nil, err
	}
	
	auth := model.Authentication{Token: tokenString, 
								ExpirationTime :expirationTime}	

	return &auth ,nil
}

func (u *UseCaseJwt) TokenValidationRSA(ctx context.Context, bearerToken string) (bool, error){
	childLogger.Debug().Msg("TokenValidationRSA")

	span := observability.Span(ctx, "useCase.TokenValidationRSA")	
    defer span.End()

	log.Debug().Interface("bearerToken : ", bearerToken).Msg("")

	claims := &model.JwtData{}

	tkn, err := jwt.ParseWithClaims(bearerToken, claims, func(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("error unexpected signing method: %v", token.Header["alg"])
		}
		return u.key_rsa_pub, nil
	})

	if err != nil {
		fmt.Println(err)
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

func (u *UseCaseJwt) RefreshTokenRSA(ctx context.Context, bearerToken string) (*model.Authentication, error){
	childLogger.Debug().Msg("RefreshTokenRSA")

	span := observability.Span(ctx, "useCase.RefreshTokenRSA")	
    defer span.End()

	// Check with token is signed 
	claims := &model.JwtData{}
	tkn, err := jwt.ParseWithClaims(bearerToken, claims, func(token *jwt.Token) (interface{}, error) {
		return u.key_rsa_pub, nil
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
	if time.Until(claims.ExpiresAt.Time) > (719 * time.Minute) {
		return nil, erro.ErrTokenStillValid
	}

	// Set a new tokens claims
	expirationTime := time.Now().Add(720 * time.Minute)
	claims.ExpiresAt = jwt.NewNumericDate(expirationTime)
	claims.ISS = "lambda-go-autentication-refreshed"

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	
	tokenString, err := token.SignedString(u.key_rsa_priv)
	if err != nil {
		return nil, err
	}

	auth := model.Authentication{	Token: tokenString, 
									ExpirationTime :expirationTime}

	return &auth,nil
}