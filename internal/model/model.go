package model

import (
	"time"
	"crypto/rsa"
	
	"github.com/golang-jwt/jwt/v4"
)

type AppServer struct {
	InfoApp 		*InfoApp 		`json:"info_app"`
	ConfigOTEL		*ConfigOTEL		`json:"otel_config"`
}

type InfoApp struct {
	AppName				string `json:"app_name,omitempty"`
	AWSRegion			string `json:"aws_region,omitempty"`
	ApiVersion			string `json:"version,omitempty"`
	TableName			string `json:"table_name,omitempty"`
	Env					string `json:"env,omitempty"`
	SecretJwtKey		string `json:"secret_jwt_key,omitempty"`
	AccountID			string `json:"account,omitempty"`
	BucketNameRSAKey	string `json:"bucket_rsa_key,omitempty"`
	FilePathRSA			string `json:"path_rsa_key,omitempty"`
	FileNameRSAPrivKey	string `json:"file_name_rsa_private_key,omitempty"`
	FileNameRSAPubKey	string `json:"file_name_rsa_public_key,omitempty"`
}

type Authentication struct {
	Token			string	`json:"token,omitempty"`
	TokenEncrypted	string	`json:"token_encrypted,omitempty"`
	ExpirationTime	time.Time `json:"expiration_time,omitempty"`
	ApiKey			string	`json:"api_key,omitempty"`
}

type Credential struct {
	ID				string	`json:"ID"`
	SK				string	`json:"SK"`
	User			string	`json:"user,omitempty"`
	Password		string	`json:"password,omitempty"`
	Token			string 	`json:"token,omitempty"`
	UsagePlan		string 	`json:"usage_plan,omitempty"`
	ApiKey			string 	`json:"apikey,omitempty"`
	Updated_at  	time.Time 	`json:"updated_at,omitempty"`
}

type CredentialScope struct {
	ID				string		`json:"ID"`
	SK				string		`json:"SK"`
	User			string		`json:"user,omitempty"`
	Scope			[]string	`json:"scope,omitempty"`
	Updated_at  	time.Time 	`json:"updated_at,omitempty"`
}

type JwtData struct {
	TokenUse	string 	`json:"token_use"`
	ISS			string 	`json:"iss"`
	Version		string 	`json:"version"`
	JwtId		string 	`json:"jwt_id"`
	Username	string 	`json:"username"`
	Scope	  	[]string `json:"scope"`
	jwt.RegisteredClaims
}

type ConfigOTEL struct {
	OtelExportEndpoint		string
	TimeInterval            int64    `mapstructure:"TimeInterval"`
	TimeAliveIncrementer    int64    `mapstructure:"RandomTimeAliveIncrementer"`
	TotalHeapSizeUpperBound int64    `mapstructure:"RandomTotalHeapSizeUpperBound"`
	ThreadsActiveUpperBound int64    `mapstructure:"RandomThreadsActiveUpperBound"`
	CpuUsageUpperBound      int64    `mapstructure:"RandomCpuUsageUpperBound"`
	SampleAppPorts          []string `mapstructure:"SampleAppPorts"`
}

type RSA_Key struct{
	SecretNameH256		string
	JwtKey				string
	Key_rsa_priv_pem	string
	Key_rsa_pub_pem 	string	
	Key_rsa_priv 		*rsa.PrivateKey
	Key_rsa_pub 		*rsa.PublicKey	
}