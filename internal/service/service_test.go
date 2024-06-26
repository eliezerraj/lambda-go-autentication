package service

import (
	"testing"
	"errors"
	"context"
	"github.com/rs/zerolog"

	"github.com/aws/aws-sdk-go-v2/config"
	
	"github.com/lambda-go-autentication/internal/core"
	"github.com/lambda-go-autentication/internal/erro"
	"github.com/lambda-go-autentication/internal/repository"

)

var (
	logLevel		= zerolog.DebugLevel // InfoLevel DebugLevel
	authService		*AuthService
	tableName		= "user_login_2"
	jwtKey			= "my_secret_key"
)

func TestSignIn(t *testing.T) {
	zerolog.SetGlobalLevel(logLevel)

	ctx := context.Background()
	awsConfig, err := config.LoadDefaultConfig(ctx)

	credential := core.Credential{User: "user123", Password: "pass123" }

	authRepository, err := repository.NewAuthRepository(tableName, awsConfig)
	if err != nil {
		t.Errorf("configuration error AuthRepository() %v ",err.Error())
	}

	authService = NewAuthService([]byte(jwtKey), authRepository)
	res, err := authService.SignIn(ctx, credential)
	if err != nil {
		t.Errorf("Error -TestSignIn Erro %v ", err)
	}
	if (res != nil) {
		t.Logf("Success TestSignIn")
	} else {
		t.Errorf("Failed TestSignIn")
	}
}

func TestAddScope(t *testing.T) {
	zerolog.SetGlobalLevel(logLevel)

	ctx := context.Background()
	awsConfig, err := config.LoadDefaultConfig(ctx)

	scope := []string{"info.read",
								"a.read",
								"sum.write",
								"version",
								"header.read"}
	credential_scope := core.CredentialScope{User: "user123", Scope: scope }

	authRepository, err := repository.NewAuthRepository(tableName, awsConfig)
	if err != nil {
		t.Errorf("configuration error AuthRepository() %v ",err.Error())
	}

	authService = NewAuthService([]byte(jwtKey), authRepository)
	res, err := authService.AddScope(ctx, credential_scope)
	if err != nil {
		t.Errorf("Error -TestAddScope Erro %v ", err)
	}
	if (res != nil) {
		t.Logf("Success TestAddScope")
	} else {
		t.Errorf("Failed TestAddScope")
	}
}

func TestQueryCredentialScope(t *testing.T) {
	zerolog.SetGlobalLevel(logLevel)

	ctx := context.Background()
	awsConfig, err := config.LoadDefaultConfig(ctx)

	credential := core.Credential{User: "user123" }

	authRepository, err := repository.NewAuthRepository(tableName, awsConfig)
	if err != nil {
		t.Errorf("configuration error AuthRepository() %v ",err.Error())
	}

	authService = NewAuthService([]byte(jwtKey), authRepository)
	res, err := authService.QueryCredentialScope(ctx, credential)
	if err != nil {
		t.Errorf("Error -TestQueryCredentialScope Erro %v ", err)
	}
	if (res != nil) {
		t.Logf("Success TestQueryCredentialScope %v ", res )
	} else {
		t.Errorf("Failed TestQueryCredentialScope")
	}
}

func TestLogin(t *testing.T) {
	zerolog.SetGlobalLevel(logLevel)

	ctx := context.Background()
	awsConfig, err := config.LoadDefaultConfig(ctx)

	credential := core.Credential{User: "user123", Password: "pass123" }

	authRepository, err := repository.NewAuthRepository(tableName, awsConfig)
	if err != nil {
		t.Errorf("configuration error AuthRepository() %v ",err.Error())
	}

	authService = NewAuthService([]byte(jwtKey), authRepository)
	res, err := authService.Login(ctx,credential)
	if err != nil {
		t.Errorf("Error -TestLogin Erro %v ", err)
	}
	if (res != nil) {
		t.Logf("Success TestLogin %v ", res)
	} else {
		t.Errorf("Failed TestLogin")
	}
}

func TestTokenValidation(t *testing.T) {
	zerolog.SetGlobalLevel(logLevel)

	ctx := context.Background()
	awsConfig, err := config.LoadDefaultConfig(ctx)
	credential := core.Credential{User: "user123", Password: "pass123" }
	
	authRepository, err := repository.NewAuthRepository(tableName, awsConfig)
	if err != nil {
		t.Errorf("configuration error AuthRepository() %v ",err.Error())
	}

	authService = NewAuthService([]byte(jwtKey), authRepository)
	res, err := authService.Login(context.TODO(),credential)
	if err != nil {
		t.Errorf("Error -TestTokenValidation Erro %v ", err)
	}
	if (res != nil) {
		t.Logf("Success TestTokenValidation/TestLogin")
	} else {
		t.Errorf("Failed TestTokenValidation/TestLogin")
	}

	credential.Token = res.Token

	isValid, err := authService.TokenValidation(ctx, credential)
	if err != nil {
		t.Errorf("Error - TestTokenValidation Erro %v ", err)
	}
	if isValid != true {
		t.Errorf("Failed - TestTokenValidation isValid %v ", isValid)
	} else {
		t.Logf("Success TestTokenValidation")
	}
}

func TestRefreshToken(t *testing.T) {
	zerolog.SetGlobalLevel(logLevel)

	ctx := context.Background()
	awsConfig, err := config.LoadDefaultConfig(ctx)
	credential := core.Credential{User: "user123", Password: "pass123" }
	
	authRepository, err := repository.NewAuthRepository(tableName, awsConfig)
	if err != nil {
		t.Errorf("configuration error AuthRepository() %v ",err.Error())
	}

	authService = NewAuthService([]byte(jwtKey), authRepository)
	res, err := authService.Login(context.TODO(),credential)
	if err != nil {
		t.Errorf("Error -TestRefreshToken Erro %v ", err)
	}
	if (res != nil) {
		t.Logf("Success TestRefreshToken/TestLogin")
	} else {
		t.Errorf("Failed TestRefreshToken/TestLogin")
	}

	credential.Token = res.Token

	_, err = authService.RefreshToken(ctx, credential)
	if errors.Is(err, erro.ErrTokenStillValid) {
		t.Logf("Success TestRefreshToken %err ", err)
	} else {
		t.Errorf("Error - TestRefreshToken Erro %v ", err)
	}
}