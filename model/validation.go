package model

import (
	"fmt"

	enum "github.com/entain-test-task/model/enum"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
)

type Validator struct {
	translator ut.Translator
	validate   *validator.Validate
}

func NewValidator() *Validator {
	validate := validator.New()

	en := en.New()
	uni := ut.New(en, en)
	translator, _ := uni.GetTranslator("en")

	enTranslations.RegisterDefaultTranslations(validate, translator)

	validate.RegisterValidation("record_state", recordStateValidator)
	validate.RegisterTranslation("record_state", translator, func(ut ut.Translator) error {
		return ut.Add("record_state", "{0} is not a valid record state", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("record_state", fe.Field())
		return t
	})

	return &Validator{
		translator: translator,
		validate:   validate,
	}
}

func (v *Validator) ValidateRequest(request interface{}) []error {
	err := v.validate.Struct(request)
	return v.translateError(err)
}

func IsArrayContainsValue(array []string, value string) bool {
	for _, item := range array {
		if item == value {
			return true
		}
	}

	return false
}

func recordStateValidator(fl validator.FieldLevel) bool {
	return IsArrayContainsValue(enum.RecordStates, fl.Field().String())
}

func (v *Validator) translateError(err error) []error {
	var errs []error

	if err == nil {
		return nil
	}

	validatorErrs := err.(validator.ValidationErrors)
	for _, e := range validatorErrs {
		translatedErr := fmt.Errorf(e.Translate(v.translator))
		errs = append(errs, translatedErr)
	}

	return errs
}
