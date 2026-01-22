package Dto

type HandlerRegister struct {
	Name     string `validate:"required,min=2,max=20"`
	Email    string `validate:"email,min=2,max=40"`
	Password string `validate:"required,gt=6,lte=25"`
}
