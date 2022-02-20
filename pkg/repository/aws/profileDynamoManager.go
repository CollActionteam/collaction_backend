package aws

import (
	"github.com/CollActionteam/collaction_backend/internal/models"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type DynamoDb struct {
	DbClient *dynamodb.DynamoDB
}

type DynamoTable struct {
	DbClient DynamoDb
	Name     string
}

type UpdateItemData struct {
	SearchKey        string
	SearchValue      string
	UpdateFieldKey   string
	UpdateFieldValue string
}

func NewDynamoConn() (svc *DynamoDb) {
	svc = &DynamoDb{DbClient: dynamodb.New(session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})))}
	return
}

func NewTable(tableName string, dbClient DynamoDb) (t *DynamoTable) {
	t = &DynamoTable{Name: tableName, DbClient: dbClient}
	return
}

func NewUpdateItem(searchKey, searchValue, updateFieldKey, updateFieldValue string) (u *UpdateItemData) {
	u = &UpdateItemData{SearchKey: searchKey, SearchValue: searchValue, UpdateFieldKey: updateFieldKey, UpdateFieldValue: updateFieldValue}
	return
}

func (t *DynamoTable) DynamoGetItemKV(key, value string, receiver interface{}) error {
	result, err := t.DbClient.DbClient.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(models.ProifleTablename),
		Key: map[string]*dynamodb.AttributeValue{
			key: {
				S: aws.String(value),
			},
		},
	})
	if err != nil {
		return err
	}

	if result.Item == nil {
		return nil
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, receiver)
	if err != nil {
		return err
	}

	return nil
}

func (t *DynamoTable) DynamoUpdateItemKV(data *UpdateItemData) error {
	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":r": {
				S: aws.String(data.UpdateFieldValue),
			},
		},
		TableName: aws.String(t.Name),
		Key: map[string]*dynamodb.AttributeValue{
			data.SearchKey: {
				S: aws.String(data.SearchValue),
			},
		},
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("set " + data.UpdateFieldKey + " = :r"),
	}

	_, err := t.DbClient.DbClient.UpdateItem(input)
	if err != nil {
		return err
	}
	return nil
}

func (t *DynamoTable) DynamoInsertItemKV(data interface{}) error {
	av, err := dynamodbattribute.MarshalMap(data)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(t.Name),
	}

	_, err = t.DbClient.DbClient.PutItem(input)
	if err != nil {
		return err
	}

	return nil
}
