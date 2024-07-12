package main

import(
	"context"

	"github.com/lambda-go-autentication/internal/service"
	"github.com/lambda-go-autentication/internal/core"
	"github.com/lambda-go-autentication/internal/util"
	"github.com/lambda-go-autentication/internal/handler"
	"github.com/lambda-go-autentication/internal/repository"
	"github.com/lambda-go-autentication/internal/config/observability"
	"github.com/lambda-go-autentication/internal/config/config_aws"
	"github.com/lambda-go-autentication/internal/config/parameter_store_aws"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/events"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
 	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda/xrayconfig"
)

var (
	logLevel		=	zerolog.DebugLevel // InfoLevel DebugLevel
	appServer		core.AppServer
	authService		*service.AuthService
	authHandler		*handler.AuthHandler
	response		*events.APIGatewayProxyResponse
	tracer 			trace.Tracer
)

func init(){
	log.Debug().Msg("init")
	zerolog.SetGlobalLevel(logLevel)
	appServer = util.GetAppInfo()
	configOTEL := util.GetOtelEnv()
	appServer.ConfigOTEL = &configOTEL
}

func main(){
	log.Debug().Msg("main")
	log.Debug().Interface("appServer :",appServer).Msg("")

	ctx := context.Background() // set config
	awsConfig, err := config_aws.GetAWSConfig(ctx)
	if err != nil {
		panic("configuration error create new aws session " + err.Error())
	}

	//----- OTEL ----//
	tp := observability.NewTracerProvider(ctx, appServer.ConfigOTEL, appServer.InfoApp)
	defer func(ctx context.Context) {
		err := tp.Shutdown(ctx)
		if err != nil {
			log.Error().Err(err).Msg("Error shutting down tracer provider")
		}
	}(ctx)

	otel.SetTextMapPropagator(xray.Propagator{})
	otel.SetTracerProvider(tp)
	tracer = tp.Tracer("lambda-go-autentication-tracer")
	//----- OTEL ----//

	clientSsm := parameter_store_aws.NewClientParameterStore(*awsConfig)
	jwtKey, err := clientSsm.GetParameter(ctx, appServer.InfoApp.SSMJwtKey)
	if err != nil {
		panic("Error GetParameter " + err.Error())
	}
	log.Debug().Str("======== > jwtKey", *jwtKey).Msg("")

	// Create a repository
	authRepository, err := repository.NewAuthRepository(appServer.InfoApp.TableName, *awsConfig)
	if err != nil {
		panic("configuration error AuthRepository(), " + err.Error())
	}
	
	authService = service.NewAuthService([]byte(*jwtKey), authRepository) // Create a authorization service and inject the repository
	authHandler = handler.NewAuthHandler(*authService, appServer) // Create a handler and inject the service

	lambda.Start(otellambda.InstrumentHandler(lambdaHandler, xrayconfig.WithRecommendedOptions(tp)... ))
}

func lambdaHandler(ctx context.Context, req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	log.Debug().Msg("lambdaHandler")
	log.Debug().Str("req.Body", req.Body).Msg("")

	ctx, span := tracer.Start(ctx, "lambdaHandler_otel_v1.2")
    defer span.End()

	// Check the http method and path
	switch req.HTTPMethod {
	case "GET":
		if (req.Resource == "/credentialScope/{id}"){  
			response, _ = authHandler.QueryCredentialScope(ctx, req) // Query the scopes associated with credential
		}else if (req.Resource == "/info"){
			response, _ = authHandler.GetInfo(ctx)
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