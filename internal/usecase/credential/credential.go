package credential

import(
	"context"
	
	"github.com/rs/zerolog/log"

	"github.com/lambda-go-autentication/pkg/observability"
	"github.com/lambda-go-autentication/internal/model"
	
	"github.com/lambda-go-autentication/internal/usecase/credential/repository"
)

var childLogger = log.With().Str("usecase", "credential").Logger()

type UseCaseCredential struct{
	repository	*repository.RepoCredential
	oAUTHToken func(context.Context, model.Credential, model.CredentialScope) (*model.Authentication, error)
	oAUTHTokenRSA func(context.Context, model.Credential, model.CredentialScope) (*model.Authentication, error)
}

func NewUseCaseCredential(	repository	*repository.RepoCredential,
							oAUTHToken func(context.Context, model.Credential, model.CredentialScope) (*model.Authentication, error),
							oAUTHTokenRSA func(context.Context, model.Credential, model.CredentialScope) (*model.Authentication, error)) *UseCaseCredential{
	childLogger.Debug().Msg("NewUseCaseCredential")

	return &UseCaseCredential{
		repository: repository,
		oAUTHToken: oAUTHToken,
		oAUTHTokenRSA: oAUTHTokenRSA,
	}
}

func (u *UseCaseCredential) SignIn(ctx context.Context, credential model.Credential) (*model.Credential, error){
	childLogger.Debug().Msg("SignIn")

	span := observability.Span(ctx, "repository.SignIn")	
    defer span.End()

	// Create a new credential
	res, err := u.repository.SignIn(ctx, credential)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (u *UseCaseCredential) Login(ctx context.Context, credential model.Credential) (*model.Authentication, error){
	childLogger.Debug().Msg("Login")
	childLogger.Debug().Interface("credential :",credential).Msg("")

	span := observability.Span(ctx, "repository.Login")	
    defer span.End()

	_, err := u.repository.Login(ctx, credential)
	if err != nil {
		childLogger.Error().Err(err).Msg("erro u.repository.Login")
		return nil, err
	}

	// get scopes associated with a credential
	credential_scope, err := u.repository.QueryCredentialScope(ctx, credential)
	if err != nil {
		childLogger.Error().Err(err).Msg("error u.repository.QueryCredentialScope")
		return nil, err
	}
	span_jwt := observability.Span(ctx, "service.create_jwt")	

	auth, err := u.oAUTHToken(ctx, credential, *credential_scope)
	if err != nil {
		childLogger.Error().Err(err).Msg("error u.oAUTHToken")
		return nil, err
	}

	defer span_jwt.End()
	return auth, nil
}

func (u *UseCaseCredential) LoginRSA(ctx context.Context, credential model.Credential) (*model.Authentication, error){
	childLogger.Debug().Msg("LoginRSA")
	childLogger.Debug().Interface("credential :",credential).Msg("")

	span := observability.Span(ctx, "repository.Login")	
    defer span.End()

	_, err := u.repository.Login(ctx, credential)
	if err != nil {
		childLogger.Error().Err(err).Msg("erro u.repository.Login")
		return nil, err
	}

	// get scopes associated with a credential
	credential_scope, err := u.repository.QueryCredentialScope(ctx, credential)
	if err != nil {
		childLogger.Error().Err(err).Msg("error u.repository.QueryCredentialScope")
		return nil, err
	}
	span_jwt := observability.Span(ctx, "service.create_jwt")	

	auth, err := u.oAUTHTokenRSA(ctx, credential, *credential_scope)
	if err != nil {
		childLogger.Error().Err(err).Msg("error u.oAUTHToken")
		return nil, err
	}

	defer span_jwt.End()
	return auth, nil
}

func (u *UseCaseCredential) AddScope(ctx context.Context, credential_scope model.CredentialScope) (*model.CredentialScope, error){
	childLogger.Debug().Msg("AddScope")

	span := observability.Span(ctx, "repository.AddScope")	
    defer span.End()

	// Save the credentials scopes
	res, err := u.repository.AddScope(ctx, credential_scope)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (u UseCaseCredential) QueryCredentialScope(ctx context.Context, credential model.Credential) (*model.CredentialScope, error){
	childLogger.Debug().Msg("QueryCredentialScope")

	span := observability.Span(ctx, "repository.QueryCredentialScope")	
    defer span.End()

	// Query all scope linked with the credentials

	childLogger.Debug().Interface("++++++++>>>>>> :", u.repository.TableName).Msg("")
	childLogger.Debug().Interface("++++++++>>>>>> :", u.repository.Repository).Msg("")

	credential_scope, err := u.repository.QueryCredentialScope(ctx, credential)
	if err != nil {
		return nil, err
	}

	return credential_scope, nil
}