package repository

import(
	"fmt"
	"time"
	"context"
	
	"github.com/rs/zerolog/log"

	"github.com/lambda-go-autentication/internal/erro"
	database "github.com/lambda-go-autentication/pkg/database/dynamo"

	"github.com/lambda-go-autentication/pkg/observability"
	"github.com/lambda-go-autentication/internal/model"
	
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var childLogger = log.With().Str("repo", "credential").Logger()

type RepoCredential struct{
	TableName   *string
	Repository	*database.Database
}

func NewRepoCredential(	repository *database.Database,
						tableName   *string) *RepoCredential{
	childLogger.Debug().Msg("NewRepoCredential")

	return &RepoCredential{
		Repository: repository,
		TableName: tableName,
	}
}

func (r *RepoCredential) SignIn(ctx context.Context, user_credential model.Credential) (*model.Credential, error){
	childLogger.Debug().Msg("SignIn")
	
	span := observability.Span(ctx, "repo.SignIn")	
    defer span.End()

	user_credential.ID 			= "USER-" + user_credential.User
	user_credential.SK 			= "USER-" + user_credential.User
	user_credential.Updated_at 	= time.Now()

	item, err := attributevalue.MarshalMap(user_credential)
	if err != nil {
		childLogger.Error().Err(err).Msg("erro MarshalMap")
		return nil, erro.ErrUnmarshal
	}

	putInput := &dynamodb.PutItemInput{
        TableName: r.TableName,
        Item:      item,
    }

	_, err = r.Repository.Client.PutItem(ctx, putInput)
    if err != nil {
		childLogger.Error().Err(err).Msg("error SignIn TransactWriteItems")
		return nil, erro.ErrInsert
    }

	return &user_credential , nil
}

func (r *RepoCredential) Login(ctx context.Context, user_credential model.Credential) (*model.Credential, error){
	childLogger.Debug().Msg("Login")

	span := observability.Span(ctx, "repo.Login")	
    defer span.End()

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

	key := &dynamodb.QueryInput{	TableName:                 r.TableName,
									ExpressionAttributeNames:  expr.Names(),
									ExpressionAttributeValues: expr.Values(),
									KeyConditionExpression:    expr.KeyCondition(),
	}

	result, err := r.Repository.Client.Query(ctx, key)
	if err != nil {
		childLogger.Error().Err(err).Msg("error Query")
		return nil, erro.ErrQuery
	}

	credential := []model.Credential{}
	err = attributevalue.UnmarshalListOfMaps(result.Items, &credential)
    if err != nil {
		childLogger.Error().Err(err).Msg("error unmarshalListOfMaps")
		return nil, erro.ErrUnmarshal
    }

	if len(credential) == 0 {
		return nil, erro.ErrNotFound
	} else {
		return &credential[0], nil
	}
}

func (r *RepoCredential) AddScope(ctx context.Context, credential_scope model.CredentialScope) (*model.CredentialScope, error){
	childLogger.Debug().Msg("AddScope")

	span := observability.Span(ctx, "repo.AddScope")	
    defer span.End()

	credential_scope.ID 			= "USER-" + credential_scope.User
	credential_scope.SK 			= "SCOPE-001"
	credential_scope.Updated_at 	= time.Now()

	item, err := attributevalue.MarshalMap(credential_scope)
	if err != nil {
		childLogger.Error().Err(err).Msg("error MarshalMap")
		return nil, erro.ErrUnmarshal
	}

	putInput := &dynamodb.PutItemInput{
        TableName: r.TableName,
        Item:      item,
    }

	_, err = r.Repository.Client.PutItem(context.TODO(), putInput)
    if err != nil {
		childLogger.Error().Err(err).Msg("error AddScope TransactWriteItems")
		return nil, erro.ErrInsert
    }

	return &credential_scope , nil
}

func (r *RepoCredential) QueryCredentialScope(ctx context.Context, user_credential model.Credential) (*model.CredentialScope, error){
	childLogger.Debug().Msg("QueryCredentialScope")

	//span := observability.Span(ctx, "repo.QueryCredentialScope")	
    //defer span.End()

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

	key := &dynamodb.QueryInput{TableName:                 r.TableName,
								ExpressionAttributeNames:  expr.Names(),
								ExpressionAttributeValues: expr.Values(),
								KeyConditionExpression:    expr.KeyCondition(),
							}

	result, err := r.Repository.Client.Query(ctx, key)
	if err != nil {
		childLogger.Error().Err(err).Msg("error Query")
		return nil, erro.ErrList
	}

	credential_scope_temp := []model.CredentialScope{}
	err = attributevalue.UnmarshalListOfMaps(result.Items, &credential_scope_temp)
    if err != nil {
		childLogger.Error().Err(err).Msg("error UnmarshalListOfMaps")
		return nil, erro.ErrUnmarshal
    }

	credential_scope_result := model.CredentialScope{}
	for _, item := range credential_scope_temp{
		credential_scope_result.ID = item.ID
		credential_scope_result.SK = item.SK
		credential_scope_result.Updated_at = item.Updated_at
		credential_scope_result.Scope = item.Scope
	}

	return &credential_scope_result, nil
}