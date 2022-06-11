package models

import (
	"context"
	"regexp"

	"github.com/CollActionteam/collaction_backend/utils"
	"github.com/go-playground/validator/v10"
)

type EmailContactRequest struct {
	Data  EmailRequestData `json:"data" validate:"required"`
	Nonce string           `json:"nonce"`
}

type EmailRequestData struct {
	Email      string `json:"email" validate:"required,email" binding:"required"`
	Subject    string `json:"subject" validate:"required,lte=50" binding:"required"`
	Message    string `json:"message" validate:"required,lte=500" binding:"required"`
	AppVersion string `json:"app_version" validate:"required" binding:"required"`
}

type EmailData struct {
	Recipient  string
	Message    string
	Subject    string
	Sender     string
	ReplyEmail string
}

func (e EmailContactRequest) Validate(ctx context.Context) validator.ValidationErrorsTranslations {
	validate := validator.New()
	if err := validate.StructCtx(ctx, e); err != nil {
		return utils.ValidationResponse(err, validate)
	}

	reg := regexp.MustCompile(`^(?:ios|android) [0-9]+\.[0-9]+\.[0-9]+\+[0-9]+$`)
	if match := reg.MatchString(e.Data.AppVersion); !match {
		return validator.ValidationErrorsTranslations{"err": "app version is not valid"}
	}

	return nil
}
