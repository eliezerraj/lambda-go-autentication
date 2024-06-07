package repository

import(
	"time"
	"fmt"
	"context"

	"github.com/lambda-go-autentication/internal/core"
	"github.com/lambda-go-autentication/internal/erro"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	
	"github.com/aws/aws-xray-sdk-go/xray"
)

func (r *AuthRepository) Login(ctx context.Context, user_credential core.Credential) (*core.Credential, error){
	childLogger.Debug().Msg("Login")

	_, root := xray.BeginSubsegment(ctx, "Repository.Login")
	defer root.Close(nil)

	var keyCond expression.KeyConditionBuilder
	id := "USER-" + user_credential.User

	keyCond = expression.KeyAnd(
		expression.Key("ID").Equal(expression.Value(id)),
		expression.Key("SK").BeginsWith(id),
	)

	expr, err := expression.NewBuilder().
							WithKeyCondition(keyCond).
							Build()
	if err != nil {
		childLogger.Error().Err(err).Msg("error NewBuilder")
		return nil, erro.ErrPreparedQuery
	}

	key := &dynamodb.QueryInput{	TableName:                 r.tableName,
									ExpressionAttributeNames:  expr.Names(),
									ExpressionAttributeValues: expr.Values(),
									KeyConditionExpression:    expr.KeyCondition(),
	}

	result, err := r.client.Query(ctx, key)
	if err != nil {
		childLogger.Error().Err(err).Msg("error Query")
		return nil, erro.ErrQuery
	}

	credential := []core.Credential{}
	err = attributevalue.UnmarshalListOfMaps(result.Items, &credential)
    if err != nil {
		childLogger.Error().Err(err).Msg("error UnmarshalListOfMaps")
		return nil, erro.ErrUnmarshal
    }

	if len(credential) == 0 {
		return nil, erro.ErrNotFound
	} else {
		return &credential[0], nil
	}
}

func (r *AuthRepository) SignIn(ctx context.Context, user_credential core.Credential) (*core.Credential, error){
	childLogger.Debug().Msg("SignIn")

	user_credential.ID 			= "USER-" + user_credential.User
	user_credential.SK 			= "USER-" + user_credential.User
	user_credential.Updated_at 	= time.Now()

	item, err := attributevalue.MarshalMap(user_credential)
	if err != nil {
		childLogger.Error().Err(err).Msg("erro MarshalMap")
		return nil, erro.ErrUnmarshal
	}

	putInput := &dynamodb.PutItemInput{
        TableName: r.tableName,
        Item:      item,
    }

	_, err = r.client.PutItem(ctx, putInput)
    if err != nil {
		childLogger.Error().Err(err).Msg("error SignIn TransactWriteItems")
		return nil, erro.ErrInsert
    }

	return &user_credential , nil
}

func (r *AuthRepository) AddScope(ctx context.Context, credential_scope core.CredentialScope) (*core.CredentialScope, error){
	childLogger.Debug().Msg("AddScope")

	credential_scope.ID 			= "USER-" + credential_scope.User
	credential_scope.SK 			= "SCOPE-001"
	credential_scope.Updated_at 	= time.Now()

	item, err := attributevalue.MarshalMap(credential_scope)
	if err != nil {
		childLogger.Error().Err(err).Msg("error MarshalMap")
		return nil, erro.ErrUnmarshal
	}

	putInput := &dynamodb.PutItemInput{
        TableName: r.tableName,
        Item:      item,
    }

	_, err = r.client.PutItem(context.TODO(), putInput)
    if err != nil {
		childLogger.Error().Err(err).Msg("error AddScope TransactWriteItems")
		return nil, erro.ErrInsert
    }

	return &credential_scope , nil
}

func (r *AuthRepository) QueryCredentialScope(ctx context.Context, user_credential core.Credential) (*core.CredentialScope, error){
	childLogger.Debug().Msg("QueryCredentialScope")

	_, root := xray.BeginSubsegment(ctx, "Repository.QueryCredentialScope")
	defer root.Close(nil)

	var keyCond expression.KeyConditionBuilder

	id := fmt.Sprintf("USER-%s", user_credential.User)
	sk := "SCOPE-001"

	keyCond = expression.KeyAnd(
		expression.Key("ID").Equal(expression.Value(id)),
		expression.Key("SK").BeginsWith(sk),
	)

	expr, err := expression.NewBuilder().
							WithKeyCondition(keyCond).
							Build()
	if err != nil {
		return nil, err
	}

	key := &dynamodb.QueryInput{TableName:                 r.tableName,
								ExpressionAttributeNames:  expr.Names(),
								ExpressionAttributeValues: expr.Values(),
								KeyConditionExpression:    expr.KeyCondition(),
							}

	result, err := r.client.Query(ctx, key)
	if err != nil {
		childLogger.Error().Err(err).Msg("error Query")
		return nil, erro.ErrList
	}

	credential_scope_temp := []core.CredentialScope{}
	err = attributevalue.UnmarshalListOfMaps(result.Items, &credential_scope_temp)
    if err != nil {
		childLogger.Error().Err(err).Msg("error UnmarshalListOfMaps")
		return nil, erro.ErrUnmarshal
    }

	credential_scope_result := core.CredentialScope{}
	for _, item := range credential_scope_temp{
		credential_scope_result.ID = item.ID
		credential_scope_result.SK = item.SK
		credential_scope_result.Updated_at = item.Updated_at
		credential_scope_result.Scope = item.Scope
	}

	return &credential_scope_result, nil
}
