package model

import (
	"fmt"

	enum "github.com/entain-test-task/model/enum"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
)

var (
	Translator ut.Translator
	Validate   *validator.Validate
)

func init() {
	validate := validator.New()

	en := en.New()
	uni := ut.New(en, en)
	Translator, _ = uni.GetTranslator("en")

	enTranslations.RegisterDefaultTranslations(validate, Translator)

	validate.RegisterValidation("record_state", recordStateValidator)
	validate.RegisterTranslation("record_state", Translator, func(ut ut.Translator) error {
		return ut.Add("record_state", "{0} is not a valid record state", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("record_state", fe.Field())
		return t
	})

	Validate = validate
}

func IsArrayContainsValue(array []string, value string) bool {
	for _, item := range array {
		if item == value {
			return true
		}
	}

	return false
}

func ValidateRequest(request interface{}) []error {
	err := Validate.Struct(request)
	return translateError(err)
}

func recordStateValidator(fl validator.FieldLevel) bool {
	return IsArrayContainsValue(enum.RecordStates, fl.Field().String())
}

func translateError(err error) (errs []error) {
	if err == nil {
		return nil
	}

	validatorErrs := err.(validator.ValidationErrors)
	for _, e := range validatorErrs {
		translatedErr := fmt.Errorf(e.Translate(Translator))
		errs = append(errs, translatedErr)
	}

	return errs
}
