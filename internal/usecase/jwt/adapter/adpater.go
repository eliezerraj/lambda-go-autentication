package adapter

import(	
	"context"
	"net/http"
	"encoding/json"

	"github.com/rs/zerolog/log"
	"github.com/aws/aws-sdk-go-v2/aws"
	
	"github.com/lambda-go-autentication/internal/usecase/jwt"

	"github.com/lambda-go-autentication/pkg/observability"
	"github.com/lambda-go-autentication/internal/model"
	"github.com/lambda-go-autentication/internal/erro"

	"github.com/aws/aws-lambda-go/events"
)

var childLogger = log.With().Str("adapter", "AdapterJwt").Logger()

type AdapterJwt struct{
	usecaseJwt	*jwt.UseCaseJwt
}

func NewAdapterJwt(usecaseJwt *jwt.UseCaseJwt) *AdapterJwt{
	childLogger.Debug().Msg("NewAdapterJwt")

	return &AdapterJwt{
		usecaseJwt: usecaseJwt,
	}
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

func (h *AdapterJwt) TokenValidation(ctx context.Context, req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	childLogger.Debug().Msg("TokenValidation")

	span := observability.Span(ctx, "adapter.tokenValidation")	
    defer span.End()

	var token model.Credential
    if err := json.Unmarshal([]byte(req.Body), &token); err != nil {
        return ApiHandlerResponse(http.StatusBadRequest, MessageBody{ErrorMsg: aws.String(err.Error())})
    }

	response, err := h.usecaseJwt.TokenValidation(ctx, token.Token)
	if err != nil {
		return ApiHandlerResponse(http.StatusUnauthorized, MessageBody{ErrorMsg: aws.String(err.Error())})
	}

	handlerResponse, err := ApiHandlerResponse(http.StatusOK, response)
	if err != nil {
		return ApiHandlerResponse(http.StatusInternalServerError, MessageBody{ErrorMsg: aws.String(err.Error())})
	}

	return handlerResponse, nil
}

func (h *AdapterJwt) TokenValidationRSA(ctx context.Context, req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	childLogger.Debug().Msg("TokenValidationRSA")

	span := observability.Span(ctx, "adapter.TokenValidationRSA")	
    defer span.End()

	var token model.Credential
    if err := json.Unmarshal([]byte(req.Body), &token); err != nil {
        return ApiHandlerResponse(http.StatusBadRequest, MessageBody{ErrorMsg: aws.String(err.Error())})
    }

	response, err := h.usecaseJwt.TokenValidationRSA(ctx, token.Token)
	if err != nil {
		return ApiHandlerResponse(http.StatusUnauthorized, MessageBody{ErrorMsg: aws.String(err.Error())})
	}

	handlerResponse, err := ApiHandlerResponse(http.StatusOK, response)
	if err != nil {
		return ApiHandlerResponse(http.StatusInternalServerError, MessageBody{ErrorMsg: aws.String(err.Error())})
	}

	return handlerResponse, nil
}

func (h *AdapterJwt) RefreshToken(ctx context.Context, req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	childLogger.Debug().Msg("RefreshToken")

	span := observability.Span(ctx, "adapter.refreshToken")	
    defer span.End()

	var token model.Credential
    if err := json.Unmarshal([]byte(req.Body), &token); err != nil {
        return ApiHandlerResponse(http.StatusBadRequest, MessageBody{ErrorMsg: aws.String(err.Error())})
    }

	response, err := h.usecaseJwt.RefreshToken(ctx, token.Token)
	if err != nil {
		return ApiHandlerResponse(http.StatusUnauthorized, MessageBody{ErrorMsg: aws.String(err.Error())})
	}

	handlerResponse, err := ApiHandlerResponse(http.StatusOK, response)
	if err != nil {
		return ApiHandlerResponse(http.StatusInternalServerError, MessageBody{ErrorMsg: aws.String(err.Error())})
	}

	return handlerResponse, nil
}

func (h *AdapterJwt) RefreshTokenRSA(ctx context.Context, req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	childLogger.Debug().Msg("RefreshTokenRSA")

	span := observability.Span(ctx, "adapter.RefreshTokenRSA")	
    defer span.End()

	var token model.Credential
    if err := json.Unmarshal([]byte(req.Body), &token); err != nil {
        return ApiHandlerResponse(http.StatusBadRequest, MessageBody{ErrorMsg: aws.String(err.Error())})
    }

	response, err := h.usecaseJwt.RefreshTokenRSA(ctx, token.Token)
	if err != nil {
		return ApiHandlerResponse(http.StatusUnauthorized, MessageBody{ErrorMsg: aws.String(err.Error())})
	}

	handlerResponse, err := ApiHandlerResponse(http.StatusOK, response)
	if err != nil {
		return ApiHandlerResponse(http.StatusInternalServerError, MessageBody{ErrorMsg: aws.String(err.Error())})
	}

	return handlerResponse, nil
}