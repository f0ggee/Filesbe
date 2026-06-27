package Dto

import (
	"Kaban/internal/DomainLevel"
	"errors"
	"github.com/go-playground/validator/v10"
)

type UserDataRegister struct {
	Name     string `validate:"required,min=2,max=20"`
	Email    string `validate:"email,min=2,max=40"`
	Password string `validate:"required,gt=6,lte=25"`
}

func (r *UserDataRegister) ValidateDate() error {

	validating := validator.New(validator.WithRequiredStructEnabled())

	err := validating.Struct(r)
	if err != nil {
		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			for _, e := range validateErrs {
				switch {
				case e.Field() == "Password":
					return r.getPasswordError(e)
				case e.Field() == "Email":
					return r.getEmailError(e)

				case e.Field() == "Name":
					return r.getNameError(e)
				}

			}
		}
	}

}

func (r *UserDataRegister) getPasswordError(e validator.FieldError) error {
	if e.Tag() == "min" {
		return errors.New(DomainLevel.PasswordSizeSmall)
	}
	if e.Tag() == "max" {
		return errors.New(DomainLevel.PasswordSizeBig)
	}
	if e.Tag() == "required" {
		return errors.New(DomainLevel.PasswordEmpty)
	}
	return errors.New(DomainLevel.NonIdentifyError)
}

func (r *UserDataRegister) getEmailError(e validator.FieldError) error {

	if e.Tag() == "email" {
		return errors.New(DomainLevel.NotCorrectEmail)
	}
	if e.Tag() == "required" {
		return errors.New(DomainLevel.EmailEmpty)
	}
	return errors.New(DomainLevel.NonIdentifyError)

}

func (r *UserDataRegister) getNameError(e validator.FieldError) error {
	if e.Tag() == "max" {
		return errors.New(DomainLevel.NameMaxSize)
	}
	if e.Tag() == "min" {
		return errors.New(DomainLevel.NameMinSize)
	}
	return errors.New(DomainLevel.NonIdentifyError)
}
