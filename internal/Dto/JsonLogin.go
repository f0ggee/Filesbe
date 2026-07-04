package Dto

import (
	"Kaban/internal/DomainLevel"
	"errors"

	"github.com/go-playground/validator/v10"
)

type UserLoginData struct {
	Email    string `validate:"email,required,min=1,max=40"`
	Password string `validate:"required,min=6,max=25"`
}

func (s *UserLoginData) ValidateData() error {
	validating := validator.New(validator.WithRequiredStructEnabled())

	err := validating.Struct(s)
	if err != nil {
		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			for _, e := range validateErrs {
				switch {
				case e.Field() == "Password":
					return s.getPasswordError(e.Field())
				case e.Field() == "Email":
					return s.getEmailError(e.Field())
				}
			}
		}
	}
	return nil
}

func (r *UserLoginData) getPasswordError(e string) error {
	if e == "min" {
		return errors.New(DomainLevel.PasswordSizeSmall)
	}
	if e == "max" {
		return errors.New(DomainLevel.PasswordSizeBig)
	}
	if e == "required" {
		return errors.New(DomainLevel.PasswordEmpty)
	}
	return errors.New(DomainLevel.NonIdentifyError)
}

func (r *UserLoginData) getEmailError(e string) error {
	if e == "email" {
		return errors.New(DomainLevel.NotCorrectEmail)
	}
	if e == "required" {
		return errors.New(DomainLevel.EmailEmpty)
	}
	return errors.New(DomainLevel.NonIdentifyError)
}
