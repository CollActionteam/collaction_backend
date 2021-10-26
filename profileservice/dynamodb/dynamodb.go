package dynamodb

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/collactionteam/collaction_backend/models"
)

var svc *dynamodb.Client

func connDb() *dynamodb.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-2"))
	if err != nil {
		log.Printf("Unable to connect to sdk config, %v", err)
	}

	svc = dynamodb.NewFromConfig(cfg)
	return svc
}

func InsertItemIntoTable(profile models.Profile, tablename string) error {
	svc = connDb()
	av, err := attributevalue.MarshalMap(profile)
	if err != nil {
		log.Println(err)
	}

	_, err = svc.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(tablename),
		Item:      av,
	})

	if err != nil {
		log.Println(err)
	}

	return nil
}
