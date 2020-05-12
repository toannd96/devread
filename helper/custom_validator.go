package helper

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/pkg/errors"
	"strings"
)

// ref > https://medium.com/@apzuk3/input-validation-in-golang-bc24cdec1835

type CustomValidator struct {
	Validator *validator.Validate
	Uni       *ut.UniversalTranslator
	Trans     ut.Translator
}

func NewCustomValidator() *CustomValidator {
	translator := en.New()
	uni := ut.New(translator, translator)
	trans, _ := uni.GetTranslator("en")

	return &CustomValidator{
		Validator: validator.New(),
		Uni:       uni,
		Trans:     trans,
	}
}

func (cv *CustomValidator) RegisterValidate() {
	if err := en_translations.RegisterDefaultTranslations(cv.Validator, cv.Trans); err != nil {
	}

	cv.Validator.RegisterValidation("pwd", func(fl validator.FieldLevel) bool {
		return len(fl.Field().String()) >= 8
	})

	cv.Validator.RegisterTranslation("required", cv.Trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0} là bắt buộc", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t
	})

	cv.Validator.RegisterTranslation("email", cv.Trans, func(ut ut.Translator) error {
		return ut.Add("email", "{0} không hợp lệ", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("email", fe.Field())
		return t
	})

	cv.Validator.RegisterTranslation("pwd", cv.Trans, func(ut ut.Translator) error {
		return ut.Add("pwd", "Mật khẩu tối thiểu 8 ký tự", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("pwd", fe.Field())
		return t
	})
}

func (cv *CustomValidator) Validate(i interface{}) error {
	err := cv.Validator.Struct(i)
	if err == nil {
		return nil
	}

	transErrors := make([]string, 0)
	for _, e := range err.(validator.ValidationErrors) {
		transErrors = append(transErrors, e.Translate(cv.Trans))
	}
	return errors.Errorf("%s", strings.Join(transErrors, " \n "))
}
