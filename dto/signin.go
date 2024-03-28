package dto

import "github.com/orewaee/go-auth/validation"

type SignInBody struct {
	Email    string `json:"email" validate:"required,email_regex"`
	Password string `json:"password" validate:"required"`
}

func (body SignInBody) Validate() error {
	return validation.GetEmailValidate().Struct(body)
}
