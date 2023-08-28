package service

import (
	"testing"
	"github.com/rs/zerolog"

	"github.com/lambda-go-autentication/internal/core/domain"

)

var (
	logLevel		=	zerolog.DebugLevel // InfoLevel DebugLevel
	authService		*AuthService
)

func TestLogin(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	jwtKey	:= "my_secret_key"
	credential := domain.Credential{User: "user123", Password: "pass123" }

	authService = NewAuthService([]byte(jwtKey))
	token, err := authService.Login(credential)
	if err != nil {
		t.Errorf("Error -TestLogin Erro %v ", err)
	}
	if (token != nil) {
		t.Logf("Success TestLogin token : %v ", token )
	} else {
		t.Errorf("Error TestLogin  %s ", credential)
	}
}

func TestTokenValidation(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.ErrorLevel)

	jwtKey	:= "my_secret_key"
	credential := domain.Credential{User: "user123", Password: "pass123" }

	authService = NewAuthService([]byte(jwtKey))
	token, err := authService.Login(credential)
	if err != nil {
		t.Errorf("Error - TestTokenValidation Erro %v ", err)
	}
	if (token == nil) {
		t.Errorf("Error TestTokenValidation Login  %s ", credential)
	}

	credential.Token = token.Token

	isValid, err := authService.TokenValidation(credential)
	if isValid != true {
		t.Errorf("Error - TestTokenValidation Erro %v ", err)
	}
}