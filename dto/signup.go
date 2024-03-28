package dto

import "github.com/orewaee/go-auth/validation"

type SignUpBody struct {
	Email    string `json:"email" validate:"required,email_regex"`
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (body SignUpBody) Validate() error {
	return validation.GetEmailValidate().Struct(body)
}
