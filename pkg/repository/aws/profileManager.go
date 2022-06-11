package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/CollActionteam/collaction_backend/internal/constants"
	"github.com/CollActionteam/collaction_backend/internal/models"
	"github.com/CollActionteam/collaction_backend/utils"
)

type Profile struct {
	dbClient Dynamo
}

func NewProfile(dynamo *Dynamo) *Profile {
	return &Profile{
		dbClient: *dynamo,
	}
}

func (p *Profile) GetUserProfile(ctx context.Context, userID string) (*models.Profile, error) {
	var profiledata *models.Profile

	err := NewTable(constants.ProfileTablename, p.dbClient).DynamoGetItemKV("userid", userID, &profiledata)
	if err != nil {
		return nil, err
	}

	profiledata.Phone, profiledata.UserID = "", ""

	return profiledata, nil
}

func (p *Profile) UpdateUserProfile(ctx context.Context, user models.UserInfo, requestData models.Profile) error {
	var (
		removeFields = []string{"userid", "displayname", "phone"}
		userID       = user.UserID
		wg           = sync.WaitGroup{}
	)

	requiredMap := map[string]string{}
	inrec, err := json.Marshal(requestData)
	if err != nil {
		return err
	}

	err = json.Unmarshal(inrec, &requiredMap)
	if err != nil {
		return err
	}

	utils.RemoveFromStringMap(requiredMap, removeFields)
	mapLength := len(requiredMap)
	if mapLength < 1 {
		return fmt.Errorf("no required update field provided")
	}

	wg.Add(mapLength)
	wrkchan := make(chan error, mapLength)
	tb := NewTable(constants.ProfileTablename, p.dbClient)

	for i, v := range requiredMap {
		go func(i string, v string, userID string, tb *DynamoTable, ch chan error, wg *sync.WaitGroup) {
			defer wg.Done()
			cData := NewUpdateItem("userid", userID, i, v)
			ch <- tb.DynamoUpdateItemKV(cData)
		}(i, v, userID, tb, wrkchan, &wg)
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

func (p *Profile) CreateUserProfile(ctx context.Context, user models.UserInfo, requestData models.Profile) error {
	var (
		profiledata *models.Profile
		tb          = NewTable(constants.ProfileTablename, p.dbClient)
	)

	err := tb.DynamoGetItemKV("userid", user.UserID, &profiledata)
	if err != nil {
		return err
	}

	if profiledata != nil {
		err := fmt.Errorf("user profile exists")
		return err
	}

	requestData.UserID, requestData.Phone = user.UserID, user.PhoneNumber
	if requestData.DisplayName == "" {
		requestData.DisplayName = user.Name
	}

	err = tb.DynamoInsertItemKV(requestData)
	if err != nil {
		return err
	}

	return nil
}
