package domain

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Authentication struct {
	Token			string	`json:"token,omitempty"`
	TokenEncrypted	string	`json:"token_encrypted,omitempty"`
	ExpirationTime	time.Time `json:"expiration_time,omitempty"`
	ApiKey			string	`json:"api_key,omitempty"`
}

type Credential struct {
	ID				string	`json:"id,omitempty"`
	SK				string	`json:"sk,omitempty"`
	User			string	`json:"user,omitempty"`
	Password		string	`json:"password,omitempty"`
	Token			string 	`json:"token,omitempty"`
	Updated_at  	time.Time 	`json:"updated_at,omitempty"`
}

type CredentialScope struct {
	ID				string		`json:"id,omitempty"`
	SK				string		`json:"sk,omitempty"`
	User			string		`json:"user,omitempty"`
	Scope			[]string	`json:"scope,omitempty"`
	Updated_at  	time.Time 	`json:"updated_at,omitempty"`
}

type JwtData struct {
	Username	string 	`json:"username"`
	Scope	  []string 	`json:"scope"`
	jwt.RegisteredClaims
}