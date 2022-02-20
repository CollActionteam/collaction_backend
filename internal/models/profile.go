package models

import (
	"os"

	validation "github.com/go-ozzo/ozzo-validation"
)

const (
	DisplayNameMinimumLength = 2
	DisplayNameMaximumLength = 20
	CountryMinimumLength     = 3
	CountryMaximumLength     = 20
	CityMinimumLength        = 3
	CityMaximumLength        = 20
	BioMinimumLength         = 10
	BioMaximumLength         = 100
)

var (
	ProifleTablename = os.Getenv("PROFILE_TABLE")
)

type Profile struct {
	UserID      string `json:"userid,omitempty"`
	DisplayName string `json:"displayname,omitempty"`
	Country     string `json:"country,omitempty"`
	City        string `json:"city,omitempty"`
	Bio         string `json:"bio,omitempty"`
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
			validation.Field(&c.DisplayName, validation.Required, validation.Length(DisplayNameMinimumLength, DisplayNameMaximumLength)),
			validation.Field(&c.Country, validation.Required, validation.Length(CountryMinimumLength, CountryMaximumLength)),
			validation.Field(&c.City, validation.Required, validation.Length(CityMinimumLength, CityMaximumLength)),
			validation.Field(&c.Bio, validation.Required, validation.Length(BioMinimumLength, BioMaximumLength)),
		)

	} else if validateType == "update" {
		return validation.ValidateStruct(&c,
			validation.Field(&c.DisplayName, validation.Length(DisplayNameMinimumLength, DisplayNameMaximumLength)),
			validation.Field(&c.Country, validation.Length(CountryMinimumLength, CountryMaximumLength)),
			validation.Field(&c.City, validation.Length(CityMinimumLength, CityMaximumLength)),
			validation.Field(&c.Bio, validation.Length(BioMinimumLength, BioMaximumLength)),
		)
	}

	return nil
}
