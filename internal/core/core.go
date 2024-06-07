package core

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type AppServer struct {
	InfoApp 		*InfoApp 		`json:"info_app"`
}

type InfoApp struct {
	AppName				string `json:"app_name,omitempty"`
	AWSRegion			string `json:"aws_region,omitempty"`
	ApiVersion			string `json:"version,omitempty"`
	AvailabilityZone 	string `json:"availabilityZone,omitempty"`
	TableName			string `json:"table_name,omitempty"`
	JwtKey				string `json:"jwt_key,omitempty"`
	SSMJwtKey			string `json:"ssm_jwt_key,omitempty"`
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
	Username	string 	`json:"username"`
	Scope	  []string 	`json:"scope"`
	jwt.RegisteredClaims
}