package main

import (
	validation "github.com/go-ozzo/ozzo-validation"
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

func (c Profile) ValidateProfileStruct(validateType string) error {
	if validateType == "create" {
		return validation.ValidateStruct(&c,
			validation.Field(&c.DisplayName, validation.Required, validation.Length(1, 20)),
			validation.Field(&c.Country, validation.Required, validation.Length(1, 20)),
			validation.Field(&c.City, validation.Required, validation.Length(1, 20)),
			validation.Field(&c.Bio, validation.Required, validation.Length(1, 100)),
		)

	} else if validateType == "update" {
		return validation.ValidateStruct(&c,
			validation.Field(&c.DisplayName, validation.Length(0, 20)),
			validation.Field(&c.Country, validation.Length(0, 20)),
			validation.Field(&c.City, validation.Length(0, 20)),
			validation.Field(&c.Bio, validation.Length(0, 100)),
		)
	}

	return nil
}
