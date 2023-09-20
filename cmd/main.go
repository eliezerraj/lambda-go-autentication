package main

import(
	"os"
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/events"

	"github.com/lambda-go-autentication/internal/service"
	"github.com/lambda-go-autentication/internal/handler"
	"github.com/lambda-go-autentication/internal/repository"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	logLevel		=	zerolog.DebugLevel // InfoLevel DebugLevel
	version			=	"1.0"
	authService		*service.AuthService
	tableName		= "user_login"
	jwtKey			= "my_secret_key"
	authHandler		*handler.AuthHandler
	response		*events.APIGatewayProxyResponse
)

// Loading ENV variables
func getEnv() {
	log.Debug().Msg("getEnv")

	if os.Getenv("LOG_LEVEL") !=  "" {
		if (os.Getenv("LOG_LEVEL") == "DEBUG"){
			logLevel = zerolog.DebugLevel
		}else if (os.Getenv("LOG_LEVEL") == "INFO"){
			logLevel = zerolog.InfoLevel
		}else if (os.Getenv("LOG_LEVEL") == "ERROR"){
				logLevel = zerolog.ErrorLevel
		}else {
			logLevel = zerolog.DebugLevel
		}
	}
	if os.Getenv("VERSION") !=  "" {
		version = os.Getenv("VERSION")
	}
	if os.Getenv("TABLE_NAME") !=  "" {
		tableName = os.Getenv("TABLE_NAME")
	}
	if os.Getenv("JWT_KEY") !=  "" {
		jwtKey = os.Getenv("JWT_KEY")
	}
}

func init(){
	log.Debug().Msg("init")
	zerolog.SetGlobalLevel(logLevel)
	getEnv()
}

func main(){
	log.Debug().Msg("main - lambda-go-autentication")

	// Create a repository
	authRepository, err := repository.NewAuthRepository(tableName)
	if err != nil {
		panic("configuration error AuthRepository(), " + err.Error())
	}
	// Create a authorization service and inject the repository
	authService = service.NewAuthService([]byte(jwtKey), authRepository)
	// Create a handler and inject the service
	authHandler = handler.NewAuthHandler(*authService)

	// Start lambda handler
	log.Debug().Msg("Start ... lambdaHandler")
	lambda.Start(lambdaHandler)
}

func lambdaHandler(ctx context.Context, req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	log.Debug().Msg("lambdaHandler")
	log.Debug().Msg("-------------------")
	log.Debug().Str("req.Body", req.Body).
				Msg("APIGateway Request.Body")
	log.Debug().Msg("--------------------")

	// Check the http method and path
	switch req.HTTPMethod {
	case "GET":
		if (req.Resource == "/credentialScope/{id}"){  
			response, _ = authHandler.QueryCredentialScope(req) // Query the scopes associated with credential
		}else {
			response, _ = authHandler.UnhandledMethod()
		}
	case "POST":
		if (req.Resource == "/login"){  
			response, _ = authHandler.Login(req) // Login
		}else if (req.Resource == "/refreshToken") {
			response, _ = authHandler.RefreshToken(req) // Refresh Token
		}else if (req.Resource == "/tokenValidation") {
			response, _ = authHandler.TokenValidation(req) // Do a JWT validation (signature and expiration date)
		}else if (req.Resource == "/signIn") {
			response, _ = authHandler.SignIn(req) // Create a new credentials
		}else if (req.Resource == "/addScope") {
			response, _ = authHandler.AddScope(req) // Add scopes to the credential
		}else {
			response, _ = authHandler.UnhandledMethod()
		}
	case "DELETE":
		response, _ = authHandler.UnhandledMethod()
	case "PUT":
		response, _ = authHandler.UnhandledMethod()
	default:
		response, _ = authHandler.UnhandledMethod()
	}

	return response, nil
}