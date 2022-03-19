package utils

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
)

func ValidationResponse(err error, validate *validator.Validate) validator.ValidationErrorsTranslations {
	errs := err.(validator.ValidationErrors)
	english := en.New()
	trans, _ := ut.New(english, english).GetTranslator("en")
	_ = enTranslations.RegisterDefaultTranslations(validate, trans)
	return errs.Translate(trans)
}
