package adapter

import(	
	"context"
	"net/http"
	"encoding/json"

	"github.com/rs/zerolog/log"
	"github.com/aws/aws-sdk-go-v2/aws"
	
	"github.com/lambda-go-autentication/internal/usecase/credential"

	"github.com/lambda-go-autentication/pkg/observability"
	"github.com/lambda-go-autentication/internal/model"
	"github.com/lambda-go-autentication/internal/erro"

	"github.com/aws/aws-lambda-go/events"
)

var childLogger = log.With().Str("adapter", "AdapterCredential").Logger()

type AdapterCredential struct{
	appServer	*model.AppServer
	useCaseCredential 	*credential.UseCaseCredential
}

func NewAdapterCredential(	appServer	*model.AppServer, 
							useCaseCredential *credential.UseCaseCredential) *AdapterCredential{
	childLogger.Debug().Msg("NewAdapterCredential")

	return &AdapterCredential{
		appServer: appServer,
		useCaseCredential: useCaseCredential,
	}
}

func (h *AdapterCredential) UnhandledMethod() (*events.APIGatewayProxyResponse, error){
	return ApiHandlerResponse(http.StatusMethodNotAllowed, MessageBody{ErrorMsg: aws.String(erro.ErrMethodNotAllowed.Error())})
}

type MessageBody struct {
	ErrorMsg 	*string `json:"error,omitempty"`
	Msg 		*string `json:"message,omitempty"`
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

func (h *AdapterCredential) SignIn(ctx context.Context, req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	childLogger.Debug().Msg("SignIn")

	span := observability.Span(ctx, "adapter.SignIn")	
    defer span.End()

	var credential model.Credential
    if err := json.Unmarshal([]byte(req.Body), &credential); err != nil {
        return ApiHandlerResponse(http.StatusBadRequest, MessageBody{ErrorMsg: aws.String(err.Error())})
    }

	response, err := h.useCaseCredential.SignIn(ctx, credential)
	if err != nil {
		return ApiHandlerResponse(http.StatusBadRequest, MessageBody{ErrorMsg: aws.String(err.Error())})
	}

	handlerResponse, err := ApiHandlerResponse(http.StatusOK, response)
	if err != nil {
		return ApiHandlerResponse(http.StatusInternalServerError, MessageBody{ErrorMsg: aws.String(err.Error())})
	}
	return handlerResponse, nil
}

func (h *AdapterCredential) Login(ctx context.Context, req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	childLogger.Debug().Msg("Login")

	span := observability.Span(ctx, "adapter.login")	
    defer span.End()

	var credential model.Credential
    if err := json.Unmarshal([]byte(req.Body), &credential); err != nil {
        return ApiHandlerResponse(http.StatusBadRequest, MessageBody{ErrorMsg: aws.String(err.Error())})
    }

	response, err := h.useCaseCredential.Login(ctx, credential)
	if err != nil {
		return ApiHandlerResponse(http.StatusNotFound, MessageBody{ErrorMsg: aws.String(err.Error())})
	}

	handlerResponse, err := ApiHandlerResponse(http.StatusOK, response)
	if err != nil {
		return ApiHandlerResponse(http.StatusInternalServerError, MessageBody{ErrorMsg: aws.String(err.Error())})
	}
	return handlerResponse, nil
}

func (h *AdapterCredential) LoginRSA(ctx context.Context, req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	childLogger.Debug().Msg("LoginRSA")

	span := observability.Span(ctx, "adapter.loginRSA")	
    defer span.End()

	var credential model.Credential
    if err := json.Unmarshal([]byte(req.Body), &credential); err != nil {
        return ApiHandlerResponse(http.StatusBadRequest, MessageBody{ErrorMsg: aws.String(err.Error())})
    }

	response, err := h.useCaseCredential.LoginRSA(ctx, credential)
	if err != nil {
		return ApiHandlerResponse(http.StatusNotFound, MessageBody{ErrorMsg: aws.String(err.Error())})
	}

	handlerResponse, err := ApiHandlerResponse(http.StatusOK, response)
	if err != nil {
		return ApiHandlerResponse(http.StatusInternalServerError, MessageBody{ErrorMsg: aws.String(err.Error())})
	}
	return handlerResponse, nil
}

func (h *AdapterCredential) AddScope(ctx context.Context, req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	childLogger.Debug().Msg("AddScope")

	span := observability.Span(ctx, "adapter.addScope")	
    defer span.End()

	var credential_scope model.CredentialScope
    if err := json.Unmarshal([]byte(req.Body), &credential_scope); err != nil {
        return ApiHandlerResponse(http.StatusBadRequest, MessageBody{ErrorMsg: aws.String(err.Error())})
    }

	response, err := h.useCaseCredential.AddScope(ctx, credential_scope)
	if err != nil {
		return ApiHandlerResponse(http.StatusBadRequest, MessageBody{ErrorMsg: aws.String(err.Error())})
	}

	handlerResponse, err := ApiHandlerResponse(http.StatusOK, response)
	if err != nil {
		return ApiHandlerResponse(http.StatusInternalServerError, MessageBody{ErrorMsg: aws.String(err.Error())})
	}
	return handlerResponse, nil
}

func (h *AdapterCredential) QueryCredentialScope(ctx context.Context, req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	childLogger.Debug().Msg("QueryCredentialScope")
	
	span := observability.Span(ctx, "adapter.QueryCredentialScope")	
    defer span.End()

	id := req.PathParameters["id"]
	if len(id) == 0 {
		return ApiHandlerResponse(http.StatusBadRequest, MessageBody{ErrorMsg: aws.String(erro.ErrQueryEmpty.Error())})
	}

	credential := model.Credential{User: id}

	response, err := h.useCaseCredential.QueryCredentialScope(ctx, credential)
	if err != nil {
		return ApiHandlerResponse(http.StatusNotFound, MessageBody{ErrorMsg: aws.String(err.Error())})
	}

	childLogger.Debug().Msg("QueryCredentialScope....4")

	handlerResponse, err := ApiHandlerResponse(http.StatusOK, response)
	if err != nil {
		return ApiHandlerResponse(http.StatusInternalServerError, MessageBody{ErrorMsg: aws.String(err.Error())})
	}
	return handlerResponse, nil
}

func (h *AdapterCredential) GetInfo(ctx context.Context) (*events.APIGatewayProxyResponse, error) {
	childLogger.Debug().Msg("GetInfo")
	
	span := observability.Span(ctx, "adapter.GetInfo")		
    defer span.End()

	handlerResponse, err := ApiHandlerResponse(http.StatusOK, h.appServer)
	if err != nil {
		return ApiHandlerResponse(http.StatusInternalServerError, MessageBody{ErrorMsg: aws.String(err.Error())})
	}

	return handlerResponse, nil
}