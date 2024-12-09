package main

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/lambda-go-autentication/pkg/handler/apigw"

	"github.com/lambda-go-autentication/internal/usecase/jwt"
	adapter_jwt"github.com/lambda-go-autentication/internal/usecase/jwt/adapter"

	"github.com/lambda-go-autentication/internal/usecase/credential"
	adapter_credential "github.com/lambda-go-autentication/internal/usecase/credential/adapter"
	"github.com/lambda-go-autentication/internal/usecase/credential/repository"

	"github.com/lambda-go-autentication/configs"
	"github.com/lambda-go-autentication/internal/model"

	"github.com/lambda-go-autentication/pkg/util"
	"github.com/lambda-go-autentication/pkg/observability"
	"github.com/lambda-go-autentication/pkg/aws_secret_manager"
	"github.com/lambda-go-autentication/pkg/aws_bucket_s3"

	database "github.com/lambda-go-autentication/pkg/database/dynamo"

	"github.com/aws/aws-lambda-go/lambda"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
 	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda/xrayconfig"
)

var (
	logLevel = zerolog.DebugLevel // InfoLevel DebugLevel
	appServer	model.AppServer
	tracer 		trace.Tracer
)

func init(){
	log.Info().Msg("init")
	zerolog.SetGlobalLevel(logLevel)

	infoApp := util.GetAppInfo()
	configOTEL := util.GetOtelEnv()

	appServer.InfoApp = &infoApp
	appServer.ConfigOTEL = &configOTEL

	log.Info().Interface("appServer : ", appServer).Msg("")
}

func main(){
	log.Info().Msg("main")

	ctx := context.Background()
	configAWS, err := configs.GetAWSConfig(ctx, appServer.InfoApp.AWSRegion)
	if err != nil {
		panic("configuration error create new aws session " + err.Error())
	}

	//Load rsa key
	clientS3 := aws_bucket_s3.NewClientS3Bucket(*configAWS)
	key_rsa_priv_pem, err := clientS3.GetObject(	ctx, 
										appServer.InfoApp.BucketNameRSAKey,
										appServer.InfoApp.FilePathRSA,
										appServer.InfoApp.FileNameRSAPrivKey)
	if err != nil {
		log.Error().Err(err).Msg("Erro GetObject")
	}
	key_rsa_pub_pem, err := clientS3.GetObject(	ctx, 
										appServer.InfoApp.BucketNameRSAKey,
										appServer.InfoApp.FilePathRSA,
										appServer.InfoApp.FileNameRSAPubKey)
	if err != nil {
		log.Error().Err(err).Msg("Erro GetObject")
	}

	//Load symetric key
	clientSecret := aws_secret_manager.NewClientSecretManager(configAWS)
	jwtKey, err := clientSecret.GetSecret(ctx, appServer.InfoApp.SecretJwtKey)
	if err != nil {
		panic("Error GetParameter, " + err.Error())
	}

	// Create client database repository
	database, err := database.NewDatabase(ctx, configAWS)
	if err != nil {
		panic("Erro repository.NewAuthRepository, " + err.Error())
	}

	// Create a usecase jwt
	useCaseJwt := jwt.NewUseCaseJwt(jwtKey, key_rsa_priv_pem, key_rsa_pub_pem)
	adapterJwt := adapter_jwt.NewAdapterJwt(useCaseJwt)

	// Create a usecase credentials
	repoCredential:= repository.NewRepoCredential(database, &appServer.InfoApp.TableName)
	useCaseCredential := credential.NewUseCaseCredential(repoCredential, useCaseJwt.OAUTHToken, useCaseJwt.OAUTHTokenRSA)
	adapterCredential := adapter_credential.NewAdapterCredential(&appServer, useCaseCredential)

	handler := apigw.InitializeLambdaHandler(adapterCredential, adapterJwt)

	tp := observability.NewTracerProvider(ctx, appServer.ConfigOTEL, appServer.InfoApp)
	defer func(ctx context.Context) {
			err := tp.Shutdown(ctx)
			if err != nil {
				log.Error().Err(err).Msg("Error shutting down tracer provider")
			}
	}(ctx)
	
	otel.SetTextMapPropagator(xray.Propagator{})
	otel.SetTracerProvider(tp)
	tracer = tp.Tracer("lambda-go-authorizer-cert")

	lambda.Start(otellambda.InstrumentHandler(handler.LambdaHandlerRequest, xrayconfig.WithRecommendedOptions(tp)... ))
	//lambda.Start(handler.LambdaHandlerRequest)
}	