package models

import (
	"github.com/CollActionteam/collaction_backend/internal/constants"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Profile struct {
	UserID      string `json:"userid,omitempty"`
	DisplayName string `json:"displayname"`
	Country     string `json:"country"`
	City        string `json:"city"`
	Bio         string `json:"bio"`
	Phone       string `json:"phone,omitempty"`
}

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Status  int         `json:"status"`
}

type UserInfo struct {
	UserID      string
	Name        string
	PhoneNumber string
}

func NewUserInfo(userId, name, phoneNumber string) *UserInfo {
	return &UserInfo{UserID: userId, Name: name, PhoneNumber: phoneNumber}
}

func (c Profile) ValidateProfileStruct(validateType string) error {
	if validateType == "create" {
		return validation.ValidateStruct(&c,
			validation.Field(&c.DisplayName, validation.Required, validation.Length(constants.DisplayNameMinimumLength, constants.DisplayNameMaximumLength)),
			validation.Field(&c.Country, validation.Required, validation.Length(constants.CountryMinimumLength, constants.CountryMaximumLength)),
			validation.Field(&c.City, validation.Required, validation.Length(constants.CityMinimumLength, constants.CityMaximumLength)),
			validation.Field(&c.Bio, validation.Required, validation.Length(constants.BioMinimumLength, constants.BioMaximumLength)),
		)

	} else if validateType == "update" {
		return validation.ValidateStruct(&c,
			validation.Field(&c.DisplayName, validation.Length(constants.DisplayNameMinimumLength, constants.DisplayNameMaximumLength)),
			validation.Field(&c.Country, validation.Length(constants.CountryMinimumLength, constants.CountryMaximumLength)),
			validation.Field(&c.City, validation.Length(constants.CityMinimumLength, constants.CityMaximumLength)),
			validation.Field(&c.Bio, validation.Length(constants.BioMinimumLength, constants.BioMaximumLength)),
		)
	}

	return nil
}
