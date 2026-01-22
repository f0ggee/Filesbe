package Helpers

import (
	"Kaban/iternal/Dto"
	"crypto/rand"
)

func CheckLenOfName(sizeAndName string) string {
	nameOfFile := sizeAndName
	if len(sizeAndName) > 5 {
		NewString := rand.Text()
		Dto.NamesToConvert[NewString[:3]] = sizeAndName
		nameOfFile = NewString[:3]

	}
	return nameOfFile
}
