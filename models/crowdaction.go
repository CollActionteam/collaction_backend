package models

import (
	"fmt"

	"github.com/CollActionteam/collaction_backend/utils"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// TODO for later: Use different model between database and api (PasswordJoin vs. IsPasswordRequired)
type Crowdaction struct {
	CrowdactionID     string             `json:"crowdactionID,omitempty"`
	Title             string             `json:"title,omitempty"`
	Description       string             `json:"description,omitempty"`
	Category          string             `json:"category,omitempty"`
	Subcategory       string             `json:"sub_category,omitempty"`
	Location          string             `json:"location,omitempty"`
	DateStart         string             `json:"date_start,omitempty"`
	DateEnd           string             `json:"date_end,omitempty"`
	DateLimitJoin     string             `json:"date_limit_join,omitempty"`
	PasswordJoin      string             `json:"password_join,omitempty"`
	CommitmentOptions []CommitmentOption `json:"commitment_options,omitempty"`
}

func GetCrowdaction(crowdactionID string, tableName string) (*Crowdaction, error) {
	var crowdaction Crowdaction
	dbClient := utils.CreateDBClient()
	val := utils.PrefixPKcrowdactionID + crowdactionID
	out, err := utils.GetDBItem(dbClient, val, tableName)
	if err != nil {
		return nil, err
	}
	if out.Item == nil {
		return nil, fmt.Errorf("crowdaction not found")
	}
	err = dynamodbattribute.UnmarshalMap(out.Item, &crowdaction)
	if err != nil {
		return nil, err
	}
	return &crowdaction, nil

}
