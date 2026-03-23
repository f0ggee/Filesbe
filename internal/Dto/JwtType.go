package Dto

import "github.com/golang-jwt/jwt/v5"

type JwtCustomStruct struct {
	UserID int
	jwt.RegisteredClaims
}
