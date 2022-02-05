package aws

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/session"
	"github.com/pkg/errors"
)

type DynamoManager struct {
	// need to look this up
	Client *dynamodb.DDB
}

func NewDynamoManager(session *session.Session) *DynamoManager {
	return &DynamoManager{Client: dynamodb.New(session)}
}

func (d *DynamoManager) GetItem(pk string) (string, error) {
	// need to look this up
	param, err := d.Client.GetItem(&dynamodb)
	if err != nil {
		return "", errors.WithStack(err)
	}

	// need to look this up
	return *param.Response, nil
}
