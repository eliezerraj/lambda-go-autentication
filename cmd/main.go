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
	"github.com/lambda-go-autentication/internal/lib"

	"go.opentelemetry.io/otel/propagation"

	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel"
 	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-sdk-go-v2/otelaws"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda/xrayconfig"

	"go.opentelemetry.io/otel/trace"
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

func InstrumentHandler(tp trace.TracerProvider, handlerFunc interface{}) interface{} {
	return otellambda.InstrumentHandler(handlerFunc,
		otellambda.WithTracerProvider(tp),
		otellambda.WithPropagator(propagation.TraceContext{}))
}

func main(){
	log.Debug().Msg("main")
	log.Debug().Interface("appServer :",appServer).Msg("")

	ctx := context.Background() // set config
	awsConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic("configuration error create new aws session " + err.Error())
	}
		
	// Instrument all AWS clients.
	otelaws.AppendMiddlewares(&awsConfig.APIOptions)

	ssmsvc := ssm.NewFromConfig(awsConfig) // Get Parameter-Store
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
	
	authService = service.NewAuthService([]byte(jwtKey), authRepository) // Create a authorization service and inject the repository
	authHandler = handler.NewAuthHandler(*authService, appServer) // Create a handler and inject the service

	//----- OTEL ----//
	tp := lib.NewTracerProvider(ctx, appServer.ConfigOTEL, appServer.InfoApp)
	defer func(ctx context.Context) {
		err := tp.Shutdown(ctx)
		if err != nil {
			log.Error().Err(err).Msg("Error shutting down tracer provider")
		}
	}(ctx)

	otel.SetTextMapPropagator(xray.Propagator{})
	otel.SetTracerProvider(tp)

	tracer = tp.Tracer("lambda-tracer-v1")
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