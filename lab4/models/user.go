package models

import "github.com/go-playground/validator/v10"

type User struct {
	Username string `validate:"required,min=3,max=32"`
	Password string `validate:"required,min=8"`
}

var validate = validator.New()

// ValidateUser validates user inputs
func ValidateUser(user User) error {
	return validate.Struct(user)
}
