package repository

import(
	"time"
	"fmt"
	"context"

	"github.com/lambda-go-autentication/internal/core/domain"
	"github.com/lambda-go-autentication/internal/erro"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-xray-sdk-go/xray"

)

func (r *AuthRepository) Login(ctx context.Context, user_credential domain.Credential) (*domain.Credential, error){
	childLogger.Debug().Msg("Login")

	_, root := xray.BeginSubsegment(ctx, "Repository.Login")
	defer root.Close(nil)

	var keyCond expression.KeyConditionBuilder
	id := "USER-" + user_credential.User

	keyCond = expression.KeyAnd(
		expression.Key("id").Equal(expression.Value(id)),
		expression.Key("sk").BeginsWith(id),
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

	result, err := r.client.QueryWithContext(ctx, key)
	if err != nil {
		childLogger.Error().Err(err).Msg("error Query")
		return nil, erro.ErrQuery
	}

	credential := []domain.Credential{}
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &credential)
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

func (r *AuthRepository) SignIn(user_credential domain.Credential) (*domain.Credential, error){
	childLogger.Debug().Msg("SignIn")

	user_credential.ID 			= "USER-" + user_credential.User
	user_credential.SK 			= "USER-" + user_credential.User
	user_credential.Updated_at 	= time.Now()

	item, err := dynamodbattribute.MarshalMap(user_credential)
	if err != nil {
		childLogger.Error().Err(err).Msg("erro MarshalMap")
		return nil, erro.ErrUnmarshal
	}

	transactItems := []*dynamodb.TransactWriteItem{}
	transactItems = append(transactItems, &dynamodb.TransactWriteItem{Put: &dynamodb.Put{
		TableName: r.tableName,
		Item:      item,
	}})

	transaction := &dynamodb.TransactWriteItemsInput{TransactItems: transactItems}
	if err := transaction.Validate(); err != nil {
		childLogger.Error().Err(err).Msg("error TransactWriteItemsInput")
		return nil, erro.ErrInsert
	}

	_, err = r.client.TransactWriteItems(transaction)
	if err != nil {
		childLogger.Error().Err(err).Msg("error TransactWriteItems")
		return nil, erro.ErrInsert
	}

	return &user_credential , nil
}

func (r *AuthRepository) AddScope(credential_scope domain.CredentialScope) (*domain.CredentialScope, error){
	childLogger.Debug().Msg("AddScope")

	credential_scope.ID 			= "USER-" + credential_scope.User
	credential_scope.SK 			= "SCOPE-001"
	credential_scope.Updated_at 	= time.Now()

	item, err := dynamodbattribute.MarshalMap(credential_scope)
	if err != nil {
		childLogger.Error().Err(err).Msg("error MarshalMap")
		return nil, erro.ErrUnmarshal
	}

	transactItems := []*dynamodb.TransactWriteItem{}
	transactItems = append(transactItems, &dynamodb.TransactWriteItem{Put: &dynamodb.Put{
		TableName: r.tableName,
		Item:      item,
	}})

	transaction := &dynamodb.TransactWriteItemsInput{TransactItems: transactItems}
	if err := transaction.Validate(); err != nil {
		childLogger.Error().Err(err).Msg("error TransactWriteItemsInput")
		return nil, erro.ErrInsert
	}

	_, err = r.client.TransactWriteItems(transaction)
	if err != nil {
		childLogger.Error().Err(err).Msg("error TransactWriteItems")
		return nil, erro.ErrInsert
	}

	return &credential_scope , nil
}

func (r *AuthRepository) QueryCredentialScope(ctx context.Context, user_credential domain.Credential) (*domain.CredentialScope, error){
	childLogger.Debug().Msg("QueryCredentialScope")
	
	_, root := xray.BeginSubsegment(ctx, "Repository.QueryCredentialScope")
	defer root.Close(nil)

	var keyCond expression.KeyConditionBuilder

	id := fmt.Sprintf("USER-%s", user_credential.User)
	sk := "SCOPE-001"

	keyCond = expression.KeyAnd(
		expression.Key("id").Equal(expression.Value(id)),
		expression.Key("sk").BeginsWith(sk),
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

	result, err := r.client.Query(key)
	if err != nil {
		childLogger.Error().Err(err).Msg("error Query")
		return nil, erro.ErrList
	}

	credential_scope_temp := []domain.CredentialScope{}
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &credential_scope_temp)
    if err != nil {
		childLogger.Error().Err(err).Msg("error UnmarshalListOfMaps")
		return nil, erro.ErrUnmarshal
    }

	credential_scope_result := domain.CredentialScope{}
	for _, item := range credential_scope_temp{
		credential_scope_result.ID = item.ID
		credential_scope_result.SK = item.SK
		credential_scope_result.Updated_at = item.Updated_at
		credential_scope_result.Scope = item.Scope
	}

	return &credential_scope_result, nil
}
