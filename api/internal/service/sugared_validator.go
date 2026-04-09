package service

import (
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
)

type SugaredErrorMessageValidator struct {
	validate *validator.Validate
}

func NewSugaredErrorMessageValidator(val *validator.Validate) *SugaredErrorMessageValidator {
	return &SugaredErrorMessageValidator{
		validate: val,
	}
}

func (s *SugaredErrorMessageValidator) TranslateValidationErrors(err error) map[string]string {
	english := en.New()
	uni := ut.New(english, english)
	trans, _ := uni.GetTranslator("en")

	enTranslations.RegisterDefaultTranslations(s.validate, trans)

	errs := err.(validator.ValidationErrors)
	errors := make(map[string]string)

	for _, e := range errs {
		errors[strings.ToLower(e.Field())] = e.Translate(trans)
	}

	return errors
}
