package main

import(
	"context"

	"github.com/lambda-go-autentication/internal/service"
	"github.com/lambda-go-autentication/internal/core"
	"github.com/lambda-go-autentication/internal/util"
	"github.com/lambda-go-autentication/internal/handler"
	"github.com/lambda-go-autentication/internal/repository"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/events"
	
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/config"
)

var (
	logLevel		=	zerolog.DebugLevel // InfoLevel DebugLevel
	appServer		core.AppServer
	authService		*service.AuthService
	authHandler		*handler.AuthHandler
	response		*events.APIGatewayProxyResponse
)

func init(){
	log.Debug().Msg("init")
	zerolog.SetGlobalLevel(logLevel)
	appServer = util.GetAppInfo()
}

func main(){
	log.Debug().Msg("main")

	// set config
	ctx := context.Background()
	awsConfig, err := config.LoadDefaultConfig(ctx)
	// Get Parameter-Store
	if err != nil {
		panic("configuration error create new aws session " + err.Error())
	}
		
	ssmsvc := ssm.NewFromConfig(awsConfig)
	param, err := ssmsvc.GetParameter(ctx, &ssm.GetParameterInput{
		Name:           aws.String(appServer.InfoApp.SSMJwtKey),
		WithDecryption: aws.Bool(false),
	})
	if err != nil {
		panic("configuration error get parameter " + err.Error())
	}
	jwtKey := *param.Parameter.Value

	log.Debug().Str("======== > ssmJwtKwy", appServer.InfoApp.SSMJwtKey).Msg("")
	log.Debug().Str("======== > jwtKey", jwtKey).Msg("")

	// Create a repository
	authRepository, err := repository.NewAuthRepository(appServer.InfoApp.TableName, awsConfig)
	if err != nil {
		panic("configuration error AuthRepository(), " + err.Error())
	}
	// Create a authorization service and inject the repository
	authService = service.NewAuthService([]byte(jwtKey), authRepository)
	// Create a handler and inject the service
	authHandler = handler.NewAuthHandler(*authService, appServer)

	// Start lambda handler
	log.Debug().Msg("Start ... lambdaHandler")
	lambda.Start(lambdaHandler)
}

func lambdaHandler(ctx context.Context, req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	log.Debug().Msg("lambdaHandler")
	//log.Debug().Msg("-------------------")
	//log.Debug().Str("req.Body", req.Body).Msg("")
	//log.Debug().Msg("--------------------")

	// Check the http method and path
	switch req.HTTPMethod {
	case "GET":
		if (req.Resource == "/credentialScope/{id}"){  
			response, _ = authHandler.QueryCredentialScope(ctx, req) // Query the scopes associated with credential
		}else if (req.Resource == "/info"){
			response, _ = authHandler.GetInfo()
		}else {
			response, _ = authHandler.UnhandledMethod()
		}
	case "POST":
		if (req.Resource == "/login"){  
			response, _ = authHandler.Login(ctx, req) // Login
		}else if (req.Resource == "/refreshToken") {
			response, _ = authHandler.RefreshToken(ctx, req) // Refresh Token
		}else if (req.Resource == "/tokenValidation") {
			response, _ = authHandler.TokenValidation(ctx, req) // Do a JWT validation (signature and expiration date)
		}else if (req.Resource == "/signIn") {
			response, _ = authHandler.SignIn(ctx, req) // Create a new credentials
		}else if (req.Resource == "/addScope") {
			response, _ = authHandler.AddScope(ctx, req) // Add scopes to the credential
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