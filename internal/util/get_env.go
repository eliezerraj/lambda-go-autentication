package util

import(
	"os"

	"github.com/rs/zerolog/log"
	"github.com/lambda-go-autentication/internal/core"
)

var childLogger = log.With().Str("util", "util").Logger()

func GetAppInfo() core.AppServer {
	childLogger.Debug().Msg("getEnv")

	var appServer	core.AppServer
	var infoApp		core.InfoApp

	if os.Getenv("APP_NAME") !=  "" {
		infoApp.AppName = os.Getenv("APP_NAME")
	}

	if os.Getenv("REGION") !=  "" {
		infoApp.AWSRegion = os.Getenv("REGION")
	}

	if os.Getenv("VERSION") !=  "" {
		infoApp.ApiVersion = os.Getenv("VERSION")
	}

	if os.Getenv("JWT_KEY") !=  "" {
		infoApp.JwtKey = os.Getenv("JWT_KEY")
	}

	if os.Getenv("SSM_JWT_KEY") !=  "" {
		infoApp.SSMJwtKey = os.Getenv("SSM_JWT_KEY")
	}

	if os.Getenv("TABLE_NAME") !=  "" {
		infoApp.TableName = os.Getenv("TABLE_NAME")
	}

	if os.Getenv("ENV") !=  "" {
		infoApp.Env = os.Getenv("ENV")
	}

	appServer.InfoApp = &infoApp

	return appServer
}