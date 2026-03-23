package Generating

import "crypto/rand"

func (g Generating) GenerateShortName() string {
	NewString := rand.Text()
	return NewString[:5]
}
