package Dto

type UserLoginData struct {
	Email    string `validate:"email,min=2,max=40"`
	Password string `validate:"required,min=6"`
}
