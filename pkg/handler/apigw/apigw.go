package apigw

import(
	"context"
	"github.com/rs/zerolog/log"

	"github.com/aws/aws-lambda-go/events"

	adapter_credential "github.com/lambda-go-autentication/internal/usecase/credential/adapter"
	adapter_jwt "github.com/lambda-go-autentication/internal/usecase/jwt/adapter"
)

var childLogger = log.With().Str("handler", "apigw").Logger()
var response		*events.APIGatewayProxyResponse

type LambdaHandler struct {
    AdapterCredential 	*adapter_credential.AdapterCredential
	AdapterJwt 			*adapter_jwt.AdapterJwt
}

func InitializeLambdaHandler( 	adapterCredential 	*adapter_credential.AdapterCredential,
								adapterJwt 			*adapter_jwt.AdapterJwt ) *LambdaHandler {
	childLogger.Debug().Msg("InitializeLambdaHandler")

    return &LambdaHandler{
        AdapterCredential: adapterCredential,
		AdapterJwt: adapterJwt,
	}
}

func (h *LambdaHandler) LambdaHandlerRequest(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	childLogger.Debug().Msg("lambdaHandlerRequest")
	
	// Check the http method and path
	switch request.HTTPMethod {
		case "GET":
			if (request.Resource == "/credentialScope/{id}"){  
				response, _ = h.AdapterCredential.QueryCredentialScope(ctx, request) // Query the scopes associated with credential
			}else if (request.Resource == "/info"){
				response, _ = h.AdapterCredential.GetInfo(ctx)
			}else {
				response, _ = h.AdapterCredential.UnhandledMethod()
			}
		case "POST":
			if (request.Resource == "/login"){  
				response, _ = h.AdapterCredential.Login(ctx, request) // Login
			}else if (request.Resource == "/loginRSA"){  
				response, _ = h.AdapterCredential.LoginRSA(ctx, request) // Login
			}else if (request.Resource == "/refreshToken") {
				response, _ = h.AdapterJwt.RefreshToken(ctx, request) // Refresh Token
			}else if (request.Resource == "/refreshTokenRSA") {
					response, _ = h.AdapterJwt.RefreshTokenRSA(ctx, request) // Refresh Token
			}else if (request.Resource == "/tokenValidation") {
				response, _ = h.AdapterJwt.TokenValidation(ctx, request) // Do a JWT validation (signature and expiration date)
			}else if (request.Resource == "/tokenValidationRSA") {
					response, _ = h.AdapterJwt.TokenValidationRSA(ctx, request) // Do a JWT validation (signature and expiration date)
			}else if (request.Resource == "/signIn") {
				response, _ = h.AdapterCredential.SignIn(ctx, request) // Create a new credentials
			}else if (request.Resource == "/addScope") {
				response, _ =  h.AdapterCredential.AddScope(ctx, request) // Add scopes to the credential
			}else {
				response, _ = h.AdapterCredential.UnhandledMethod()
			}
		case "DELETE":
			response, _ = h.AdapterCredential.UnhandledMethod()
		case "PUT":
			response, _ = h.AdapterCredential.UnhandledMethod()
		default:
			response, _ = h.AdapterCredential.UnhandledMethod()
	}

	return response, nil
}