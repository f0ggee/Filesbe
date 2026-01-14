package Dto

import "github.com/golang-jwt/jwt/v5"

type MyCustomCookie struct {
	UserID int
	Sc     string
	jwt.RegisteredClaims
}
