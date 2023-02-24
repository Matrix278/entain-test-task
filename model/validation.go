package model

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
)

var Validate *validator.Validate

func init() {
	validate := validator.New()

	english := en.New()
	uni := ut.New(english, english)
	trans, _ := uni.GetTranslator("en")
	_ = enTranslations.RegisterDefaultTranslations(validate, trans)

	validate.RegisterValidation("record_state", RecordStateValidator)

	Validate = validate
}

func RecordStateValidator(fl validator.FieldLevel) bool {
	return IsArrayContainsValue(recordStateValues, fl.Field().String())
}

func IsArrayContainsValue(array []string, value string) bool {
	for _, item := range array {
		if item == value {
			return true
		}
	}

	return false
}

func ValidateRequest(request interface{}) error {
	return Validate.Struct(request)
}
