package Uttiltesss

import "golang.org/x/crypto/bcrypt"

func HashPassowrd(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err

}
