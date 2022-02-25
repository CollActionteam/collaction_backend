package aws

import (
	"fmt"

	"github.com/CollActionteam/collaction_backend/internal/constants"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const (
	PartitionKey = "pk"
	SortKey      = "sk"

	CrowdactionsPageLength = 50
	PKCrowdaction          = "act"

	//access pattern getParticipation
	//item has PK="prt#usr#<userID>" and SK="prt#act#<crowdactionID>"
	//(we want strong consistency when listing the users participation)
	prefixParticipationKey              = "prt#"
	PrefixParticipationPK_UserID        = prefixParticipationKey + "usr#"
	PrefixParticipationSK_CrowdactionID = prefixParticipationKey + PKCrowdaction + "#"
)

type PrimaryKey map[string]*dynamodb.AttributeValue

type Dynamo struct {
	dbClient *dynamodb.DynamoDB
}

type DynamoTable struct {
	DbClient Dynamo
	Name     string
}

type UpdateItemData struct {
	SearchKey        string
	SearchValue      string
	UpdateFieldKey   string
	UpdateFieldValue string
}

func NewDynamo() *Dynamo {
	sess := session.Must(session.NewSession())
	return &Dynamo{dbClient: dynamodb.New(sess)}
}

func NewTable(tableName string, dbClient Dynamo) (t *DynamoTable) {
	t = &DynamoTable{Name: tableName, DbClient: dbClient}
	return
}

func NewUpdateItem(searchKey, searchValue, updateFieldKey, updateFieldValue string) (u *UpdateItemData) {
	u = &UpdateItemData{SearchKey: searchKey, SearchValue: searchValue, UpdateFieldKey: updateFieldKey, UpdateFieldValue: updateFieldValue}
	return
}

func (s *Dynamo) GetPrimaryKey(pk string, sk string) PrimaryKey {
	return PrimaryKey{
		PartitionKey: {
			S: aws.String(pk),
		},
		SortKey: {
			S: aws.String(sk),
		},
	}
}

func (s *Dynamo) GetDBItem(tableName string, pk string, sk string) (map[string]*dynamodb.AttributeValue, error) {
	result, err := s.dbClient.GetItem(&dynamodb.GetItemInput{
		TableName: &tableName,
		Key:       s.GetPrimaryKey(pk, sk),
	})

	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == dynamodb.ErrCodeResourceNotFoundException {
				err = nil // Just return nil (not found is not an error)
			}
		}
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	return result.Item, nil
}

func (s *Dynamo) PutDBItem(tableName string, pk string, sk string, record interface{}) error {
	av, err := dynamodbattribute.MarshalMap(record)
	if err != nil {
		return err
	}
	if _, hasKey := av[PartitionKey]; hasKey {
		return fmt.Errorf("record must not have a field with the label \"pk\"")
	}
	if _, hasKey := av[SortKey]; hasKey {
		return fmt.Errorf("record must not have a field with the label \"sk\"")
	}
	av[PartitionKey] = &dynamodb.AttributeValue{S: aws.String(pk)}
	av[SortKey] = &dynamodb.AttributeValue{S: aws.String(sk)}
	_, err = s.dbClient.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      av,
	})
	return err
}

func (s *Dynamo) DeleteDBItem(tableName string, pk string, sk string) error {
	_, err := s.dbClient.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: &tableName,
		Key:       s.GetPrimaryKey(pk, sk),
	})
	return err
}

func (t *DynamoTable) DynamoGetItemKV(key, value string, receiver interface{}) error {
	result, err := t.DbClient.dbClient.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(constants.ProfileTablename),
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

	_, err := t.DbClient.dbClient.UpdateItem(input)
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

	_, err = t.DbClient.dbClient.PutItem(input)
	if err != nil {
		return err
	}

	return nil
}
