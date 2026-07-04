package Dto

import (
	"Kaban/internal/DomainLevel"
	"errors"

	"github.com/go-playground/validator/v10"
)

type UserDataRegister struct {
	Name     string `validate:"required,min=2,max=20"`
	Email    string `validate:"email,required,min=1,max=40"`
	Password string `validate:"required,min=6,max=25"`
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
					return r.getPasswordError(e.Tag())
				case e.Field() == "Email":
					return r.getEmailError(e.Tag())
				case e.Field() == "Name":
					return r.getNameError(e.Tag())
				}
			}
		}
	}
	return nil
}

func (r *UserDataRegister) getPasswordError(e string) error {
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

func (r *UserDataRegister) getEmailError(e string) error {
	if e == "email" {
		return errors.New(DomainLevel.NotCorrectEmail)
	}
	if e == "min" {
		return errors.New(DomainLevel.NotCorrectEmail)
	}
	if e == "required" {
		return errors.New(DomainLevel.EmailEmpty)
	}
	if e == "max" {
		return errors.New(DomainLevel.EmailMaxSize)
	}
	return errors.New(DomainLevel.NonIdentifyError)
}

func (r *UserDataRegister) getNameError(e string) error {
	if e == "max" {
		return errors.New(DomainLevel.NameMaxSize)
	}
	if e == "min" {
		return errors.New(DomainLevel.NameMinSize)
	}
	return errors.New(DomainLevel.NonIdentifyError)
}
