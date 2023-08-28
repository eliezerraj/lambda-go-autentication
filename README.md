# lambda-go-autentication

POC Lambda for technical purposes

Lambda mock a login and return a JWT Oath

## Compile lambda

   Manually compile the function

      GOOD=linux GOARCH=amd64 go build -o ../build/main main.go

      zip -jrm ../build/main.zip ../build/main

      aws lambda update-function-code \
        --function-name lambda-go-autentication \
        --zip-file fileb:///mnt/c/Eliezer/workspace/github.com/lambda-go-autentication/build/main.zip \
        --publish

## Endpoints

+ POST /login

      {
         "User": "007",
         "Password": "MrBeam",
      }


+ POST /tokenvalidation

      {
         "Token": "ABC123",
      }

## Pipeline

Prerequisite: 

Lambda function already created

+ buildspec.yml: build the main.go and move to S3
+ buildspec-test.yml: make a go test using services_test.go
+ buildspec-update.yml: update the lambda-function using S3 build