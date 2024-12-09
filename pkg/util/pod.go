package util

import(
	"os"

	"github.com/rs/zerolog/log"
	"github.com/lambda-go-autentication/internal/model"
)

func GetAppInfo() model.InfoApp {
	log.Debug().Msg("GetAppInfo")

	var infoApp		model.InfoApp

	if os.Getenv("APP_NAME") !=  "" {
		infoApp.AppName = os.Getenv("APP_NAME")
	}

	if os.Getenv("REGION") !=  "" {
		infoApp.AWSRegion = os.Getenv("REGION")
	}

	if os.Getenv("VERSION") !=  "" {
		infoApp.ApiVersion = os.Getenv("VERSION")
	}

	if os.Getenv("SECRET_JWT_KEY") !=  "" {
		infoApp.SecretJwtKey = os.Getenv("SECRET_JWT_KEY")
	}

	if os.Getenv("TABLE_NAME") !=  "" {
		infoApp.TableName = os.Getenv("TABLE_NAME")
	}

	if os.Getenv("RSA_BUCKET_NAME_KEY") !=  "" {
		infoApp.BucketNameRSAKey = os.Getenv("RSA_BUCKET_NAME_KEY")
	}

	if os.Getenv("RSA_FILE_PATH") !=  "" {
		infoApp.FilePathRSA = os.Getenv("RSA_FILE_PATH")
	}

	if os.Getenv("RSA_PRIV_FILE_KEY") !=  "" {
		infoApp.FileNameRSAPrivKey = os.Getenv("RSA_PRIV_FILE_KEY")
	}

	if os.Getenv("RSA_PUB_FILE_KEY") !=  "" {
		infoApp.FileNameRSAPubKey = os.Getenv("RSA_PUB_FILE_KEY")
	}

	return infoApp
}