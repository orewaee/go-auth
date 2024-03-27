package validation

import (
	"github.com/go-playground/validator/v10"
	"log"
	"regexp"
)

var validate = validator.New()

func init() {
	err := validate.RegisterValidation("email_regex", func(field validator.FieldLevel) bool {
		rxp, _ := regexp.Compile(`.*@.*[.].*`)

		return rxp.MatchString(field.Field().String())
	})

	if err != nil {
		log.Fatalln(err)
	}
}

func GetEmailValidate() *validator.Validate {
	return validate
}
