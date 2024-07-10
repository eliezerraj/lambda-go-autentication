package handler

import(
	"net/http"
	"encoding/json"
	"context"

	"github.com/rs/zerolog/log"
	"github.com/lambda-go-autentication/internal/service"
	
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-lambda-go/events"

	"github.com/lambda-go-autentication/internal/config/observability"
	"github.com/lambda-go-autentication/internal/erro"
	"github.com/lambda-go-autentication/internal/core"
)

var childLogger = log.With().Str("handler", "handler").Logger()

type AuthHandler struct {
	authService service.AuthService
	appServer		core.AppServer
}

type MessageBody struct {
	ErrorMsg 	*string `json:"error,omitempty"`
	Msg 		*string `json:"message,omitempty"`
}

func NewAuthHandler(authService service.AuthService,
					appServer	core.AppServer) *AuthHandler{
	childLogger.Debug().Msg("NewAuthHandler")
	return &AuthHandler{
		authService: authService,
		appServer:	appServer,
	}
}

func ApiHandlerResponse(statusCode int, body interface{}) (*events.APIGatewayProxyResponse, error){
	stringBody, err := json.Marshal(&body)
	if err != nil {
		return nil, erro.ErrUnmarshal
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(stringBody),
	}, nil
}

func (h *AuthHandler) UnhandledMethod() (*events.APIGatewayProxyResponse, error){
	return ApiHandlerResponse(http.StatusMethodNotAllowed, MessageBody{ErrorMsg: aws.String(erro.ErrMethodNotAllowed.Error())})
}

func (h *AuthHandler) Login(ctx context.Context, req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	childLogger.Debug().Msg("Login")

	span := observability.Span(ctx, "handler.login")	
    defer span.End()

	var credential core.Credential
    if err := json.Unmarshal([]byte(req.Body), &credential); err != nil {
        return ApiHandlerResponse(http.StatusBadRequest, MessageBody{ErrorMsg: aws.String(err.Error())})
    }

	response, err := h.authService.Login(ctx, credential)
	if err != nil {
		return ApiHandlerResponse(http.StatusNotFound, MessageBody{ErrorMsg: aws.String(err.Error())})
	}

	handlerResponse, err := ApiHandlerResponse(http.StatusOK, response)
	if err != nil {
		return ApiHandlerResponse(http.StatusInternalServerError, MessageBody{ErrorMsg: aws.String(err.Error())})
	}
	return handlerResponse, nil
}

func (h *AuthHandler) SignIn(ctx context.Context, req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	childLogger.Debug().Msg("SignIn")

	span := observability.Span(ctx, "handler.signIn")	
    defer span.End()

	var credential core.Credential
    if err := json.Unmarshal([]byte(req.Body), &credential); err != nil {
        return ApiHandlerResponse(http.StatusBadRequest, MessageBody{ErrorMsg: aws.String(err.Error())})
    }

	response, err := h.authService.SignIn(ctx, credential)
	if err != nil {
		return ApiHandlerResponse(http.StatusBadRequest, MessageBody{ErrorMsg: aws.String(err.Error())})
	}

	handlerResponse, err := ApiHandlerResponse(http.StatusOK, response)
	if err != nil {
		return ApiHandlerResponse(http.StatusInternalServerError, MessageBody{ErrorMsg: aws.String(err.Error())})
	}
	return handlerResponse, nil
}

func (h *AuthHandler) TokenValidation(ctx context.Context, req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	childLogger.Debug().Msg("TokenValidation")

	span := observability.Span(ctx, "handler.tokenValidation")	
    defer span.End()

	var token core.Credential
    if err := json.Unmarshal([]byte(req.Body), &token); err != nil {
        return ApiHandlerResponse(http.StatusBadRequest, MessageBody{ErrorMsg: aws.String(err.Error())})
    }

	response, err := h.authService.TokenValidation(ctx, token)
	if err != nil {
		return ApiHandlerResponse(http.StatusUnauthorized, MessageBody{ErrorMsg: aws.String(err.Error())})
	}

	handlerResponse, err := ApiHandlerResponse(http.StatusOK, response)
	if err != nil {
		return ApiHandlerResponse(http.StatusInternalServerError, MessageBody{ErrorMsg: aws.String(err.Error())})
	}

	return handlerResponse, nil
}

func (h *AuthHandler) RefreshToken(ctx context.Context, req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	childLogger.Debug().Msg("RefreshToken")

	span := observability.Span(ctx, "handler.refreshToken")	
    defer span.End()

	var token core.Credential
    if err := json.Unmarshal([]byte(req.Body), &token); err != nil {
        return ApiHandlerResponse(http.StatusBadRequest, MessageBody{ErrorMsg: aws.String(err.Error())})
    }

	response, err := h.authService.RefreshToken(ctx, token)
	if err != nil {
		return ApiHandlerResponse(http.StatusUnauthorized, MessageBody{ErrorMsg: aws.String(err.Error())})
	}

	handlerResponse, err := ApiHandlerResponse(http.StatusOK, response)
	if err != nil {
		return ApiHandlerResponse(http.StatusInternalServerError, MessageBody{ErrorMsg: aws.String(err.Error())})
	}

	return handlerResponse, nil
}

func (h *AuthHandler) QueryCredentialScope(ctx context.Context, req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	childLogger.Debug().Msg("QueryCredentialScope")
	
	span := observability.Span(ctx, "handler.queryCredentialScope")	
    defer span.End()

	id := req.PathParameters["id"]
	if len(id) == 0 {
		return ApiHandlerResponse(http.StatusBadRequest, MessageBody{ErrorMsg: aws.String(erro.ErrQueryEmpty.Error())})
	}

	credential := core.Credential{User: id}
	response, err := h.authService.QueryCredentialScope(ctx , credential)
	if err != nil {
		return ApiHandlerResponse(http.StatusNotFound, MessageBody{ErrorMsg: aws.String(err.Error())})
	}

	handlerResponse, err := ApiHandlerResponse(http.StatusOK, response)
	if err != nil {
		return ApiHandlerResponse(http.StatusInternalServerError, MessageBody{ErrorMsg: aws.String(err.Error())})
	}
	return handlerResponse, nil
}

func (h *AuthHandler) AddScope(ctx context.Context, req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	childLogger.Debug().Msg("AddScope")

	span := observability.Span(ctx, "handler.addScope")	
    defer span.End()

	var credential_scope core.CredentialScope
    if err := json.Unmarshal([]byte(req.Body), &credential_scope); err != nil {
        return ApiHandlerResponse(http.StatusBadRequest, MessageBody{ErrorMsg: aws.String(err.Error())})
    }

	response, err := h.authService.AddScope(ctx, credential_scope)
	if err != nil {
		return ApiHandlerResponse(http.StatusBadRequest, MessageBody{ErrorMsg: aws.String(err.Error())})
	}

	handlerResponse, err := ApiHandlerResponse(http.StatusOK, response)
	if err != nil {
		return ApiHandlerResponse(http.StatusInternalServerError, MessageBody{ErrorMsg: aws.String(err.Error())})
	}
	return handlerResponse, nil
}

func (h *AuthHandler) GetInfo(ctx context.Context) (*events.APIGatewayProxyResponse, error) {
	childLogger.Debug().Msg("GetInfo")
	
	span := observability.Span(ctx, "handler.getInfo")	
    defer span.End()

	handlerResponse, err := ApiHandlerResponse(http.StatusOK, h.appServer)
	if err != nil {
		return ApiHandlerResponse(http.StatusInternalServerError, MessageBody{ErrorMsg: aws.String(err.Error())})
	}

	return handlerResponse, nil
}