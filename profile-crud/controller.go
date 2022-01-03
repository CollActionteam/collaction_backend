package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/CollActionteam/collaction_backend/auth"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
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

func UpdateProfile(req events.APIGatewayV2HTTPRequest) (err error) {
	var profiledata = Profile{}
	var skipFields = []string{"userid", "displayname", "phone"}

	svc := connDb()

	usrInf, err := auth.ExtractUserInfo(req)
	if err != nil {
		return err
	}

	userID := usrInf.UserID()

	err = json.Unmarshal([]byte(req.Body), &profiledata)
	if err != nil {
		return err
	}

	err = profiledata.ValidateProfileStruct("update")
	if err != nil {
		return err
	}

	requiredMap := map[string]string{}

	inrec, err := json.Marshal(profiledata)
	if err != nil {
		return err
	}

	err = json.Unmarshal(inrec, &requiredMap)
	if err != nil {
		return err
	}

	DeleteMapProps(requiredMap, skipFields)

	mapLength := len(requiredMap)

	if mapLength < 1 {
		return fmt.Errorf("no required update field provided")
	}

	var wg sync.WaitGroup

	wg.Add(mapLength)

	wrkchan := make(chan error, mapLength)

	for i, v := range requiredMap {
		go UpdateProfileTableItem(i, v, userID, svc, wrkchan, &wg)
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

func UpdateProfileTableItem(i string, v string, userID string, svc *dynamodb.DynamoDB, ch chan error, wg *sync.WaitGroup) {
	defer wg.Done()

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":r": {
				S: aws.String(v),
			},
		},
		TableName: aws.String(tablename),
		Key: map[string]*dynamodb.AttributeValue{
			"userid": {
				S: aws.String(userID),
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

func GetProfile(req events.APIGatewayV2HTTPRequest) (*Profile, error) {
	var profiledata *Profile
	var userID string

	svc := connDb()

	idParameter := req.PathParameters["userID"]
	if idParameter != "" {
		userID = idParameter
	} else {
		return nil, errors.New("no profile selected")
	}

	searchResult, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tablename),
		Key: map[string]*dynamodb.AttributeValue{
			"userid": {
				S: aws.String(userID),
			},
		},
	})
	if err != nil {
		return nil, err
	}

	if searchResult.Item == nil {
		return nil, nil
	}

	err = dynamodbattribute.UnmarshalMap(searchResult.Item, &profiledata)
	if err != nil {
		return nil, err
	}

	profiledata.Phone, profiledata.UserID = "", ""

	return profiledata, nil
}

func CreateProfile(req events.APIGatewayV2HTTPRequest) error {
	var profiledata Profile

	svc := connDb()

	usrInf, err := auth.ExtractUserInfo(req)
	if err != nil {
		return err
	}

	userID := usrInf.UserID()
	userName := usrInf.Name()
	userPhoneNumber := usrInf.PhoneNumber()

	searchResult, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tablename),
		Key: map[string]*dynamodb.AttributeValue{
			"userid": {
				S: aws.String(userID),
			},
		},
	})
	if err != nil {
		return err
	}

	if searchResult.Item != nil {
		err := fmt.Errorf("user profile exists")
		return err
	}

	err = json.Unmarshal([]byte(req.Body), &profiledata)
	if err != nil {
		return err
	}

	err = profiledata.ValidateProfileStruct("create")
	if err != nil {
		return err
	}

	profiledata.UserID = userID
	profiledata.Phone = userPhoneNumber

	if profiledata.DisplayName == "" {
		profiledata.DisplayName = userName
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

func DeleteMapProps(m map[string]string, s []string) {
	for _, v := range s {
		delete(m, v)
	}
}
