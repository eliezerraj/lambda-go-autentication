# lambda-go-autentication

POC Lambda for technical purposes

Lambda mock a login and return a JWT/Scope Oath using a HS256 (symetric key) The JWT token is 60 minutes duration

It saves the credentials and scopes in a DynamoDB table

![Alt text](/assets/image.png)

See: lambda-go-auth-apigw (extend example)

## Enviroment variable

+ tablename: DynamoDB table

+ jwtKey: The KEY used for encrypt Hs256

## Compile lambda

   Manually compile the function

      Old Version 
      GOOD=linux GOARCH=amd64 go build -o ../build/main main.go
      zip -jrm ../build/main.zip ../build/main

      Convert
      aws lambda update-function-configuration --function-name lambda-go-autentication --runtime provided.al2

      New Version
      GOARCH=amd64 GOOS=linux go build -o ../build/bootstrap main.go
      zip -jrm ../build/main.zip ../build/bootstrap

      aws lambda update-function-code \
        --region us-east-2 \
        --function-name lambda-go-autentication \
        --zip-file fileb:///mnt/c/Eliezer/workspace/github.com/lambda-go-autentication/build/main.zip \
        --publish

## Install a LambdaLayer (OTEL)

arn:aws:lambda:us-east-2:901920570463:layer:aws-otel-collector-amd64-ver-0-90-1:1

## Endpoints

+ POST /signIn

      {
         "user": "007",
         "password": "MrBeam",
      }

+ POST /login

      {
         "user": "007",
         "password": "MrBeam",
      }

+ POST /tokenValidation

      {
         "token": "ABC123",
      }

+ POST /refreshToken

      {
         "token": "ABC123",
      }

+ POST /addScope

      {
         "user": "user-01",
         "scope": ["test.read","test.write"]
      }

      or

      {
         "user": "user-01",
         "scope": ["admin"]
      }

      or

      {
         "user": "user-01",
         "scope": ["info"]
      }

+ GET /credentialScope/user-01

      {
         "id": "USER-user-02",
         "sk": "SCOPE-001",
         "scope": [
            "header.read",
            "version.read",
            "info.read"
         ],
         "updated_at": "2023-09-11T01:29:54.7366791Z"
      }

## Pipeline

Prerequisite: 

Lambda function already created

+ buildspec.yml: build the main.go and move to S3
+ buildspec-test.yml: make a go test using services_test.go
+ buildspec-update.yml: update the lambda-function using S3 build
+ appspec.yml: blue/gree deploy

## Lambda Env Variables

      ENV:        dev  
      APP_NAME:	lambda-go-autentication
      JWT_KEY:	   my_secret_key
      REGION:	   us-east-2
      SSM_JWT_KEY:	key-secret
      TABLE_NAME:	   user_login_2
      OTEL_EXPORTER_OTLP_ENDPOINT: localhost:4317

## Running locally

+ Create a docker image

      docker build -t lambda-go-autentication .

+ Setup the docker compose
+ Download the lambda aws-lambda-rie

      mkdir -p .aws-lambda-rie && curl -Lo .aws-lambda-rie/aws-lambda-rie https://github.com/aws/aws-lambda-runtime-interface-emulator/releases/latest/download/aws-lambda-rie && chmod +x .aws-lambda-rie/aws-lambda-rie

+ start the docker compose
   /deployment-locally/start.sh

+ Test

      curl --location 'http://localhost:9000/2015-03-31/functions/function/invocations' --header 'Content-Type: application/json' --data '{"httpMethod":"GET","resource":"/info","pathParameters": {"id":"1"}}'