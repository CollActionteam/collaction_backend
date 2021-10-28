package profileservice

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/CollActionteam/collaction_backend/auth"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/mitchellh/mapstructure"
)

var (
	tablename = os.Getenv("PROFILE_TABLE")
)

func connDb() (svc *dynamodb.DynamoDB) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc = dynamodb.New(sess)
	return
}

func UpdateProfile(req events.APIGatewayProxyRequest) (err error) {
	var profiledata = Profile{}

	svc := connDb()

	usrInf, err := auth.ExtractUserInfo(req)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(req.Body), &profiledata)
	if err != nil {
		return err
	}

	requiredMap := map[string]*dynamodb.AttributeValue{}

	err = mapstructure.Decode(profiledata, &requiredMap)
	if err != nil {
		return err
	}

	nw := len(requiredMap)

	if nw < 1 {
		return fmt.Errorf("No user data provided")
	}

	var wg sync.WaitGroup

	wg.Add(nw)

	wrkchan := make(chan error, nw)

	for i, v := range requiredMap {
		varMap := map[string]*dynamodb.AttributeValue{}
		varMap[i] = v
		go UpdateProfileTableItem(i, v, usrInf.UserID(), svc, wrkchan, &wg)
	}

	go func() {
		defer close(wrkchan)
		wg.Wait()
	}()

	for n := range wrkchan {
		if n != nil {
			err = n
		}
	}

	if err != nil {
		return err
	}

	return nil
}

func UpdateProfileTableItem(i string, v *dynamodb.AttributeValue, userID string, svc *dynamodb.DynamoDB, ch chan error, wg *sync.WaitGroup) {
	defer wg.Done()

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":r": {
				N: aws.String(fmt.Sprintf("%v", v)),
			},
		},
		TableName: aws.String(tablename),
		Key: map[string]*dynamodb.AttributeValue{
			"UserID": {
				N: aws.String(userID),
			},
		},
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("set " + i + " = :r"),
	}

	_, err := svc.UpdateItem(input)
	if err != nil {
		ch <- err
		return
	}

	ch <- nil
}

func GetProfile(req events.APIGatewayProxyRequest) (Profile, error) {
	var profiledata = Profile{}

	svc := connDb()

	// usrInf, err := auth.ExtractUserInfo(req)
	// if err != nil {
	// 	return Profile{}, err
	// }

	// userID := usrInf.UserID()
	userID := "WHHl9feFyeXyU0QwwASGxqsFgTN2"

	searchResult, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tablename),
		Key: map[string]*dynamodb.AttributeValue{
			"UserID": {
				N: aws.String(userID),
			},
		},
	})
	if err != nil {
		return Profile{}, err
	}

	if searchResult.Item == nil {
		return Profile{}, fmt.Errorf("user profile does not exist")
	}

	err = dynamodbattribute.UnmarshalMap(searchResult.Item, &profiledata)
	if err != nil {
		return Profile{}, err
	}

	return profiledata, nil
}

func CreateProfile(req events.APIGatewayProxyRequest) error {
	var profiledata Profile

	svc := connDb()

	usrInf, err := auth.ExtractUserInfo(req)
	if err != nil {
		return err
	}

	searchResult, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tablename),
		Key: map[string]*dynamodb.AttributeValue{
			"UserID": {
				N: aws.String(usrInf.UserID()),
			},
		},
	})
	if err != nil {
		return err
	}

	if searchResult.Item != nil {
		return fmt.Errorf("user profile exists")
	}

	err = json.Unmarshal([]byte(req.Body), &profiledata)
	if err != nil {
		return err
	}

	profiledata.UserID = usrInf.UserID()

	if profiledata.DisplayName == "" {
		profiledata.DisplayName = usrInf.Name()
	}

	if profiledata.Phone == "" {
		profiledata.DisplayName = usrInf.PhoneNumber()
	}

	err = InsertItemIntoProfileTable(profiledata, svc)
	if err != nil {
		return err
	}

	return nil
}

func InsertItemIntoProfileTable(profile Profile, svc *dynamodb.DynamoDB) error {
	av, err := dynamodbattribute.MarshalMap(profile)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tablename),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		return err
	}

	return nil
}
