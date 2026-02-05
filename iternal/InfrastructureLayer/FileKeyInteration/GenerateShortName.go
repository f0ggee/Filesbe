package FileKeyInteration

import "crypto/rand"

func (*FileInfoController) GenerateShortFileName() string {
	NewString := rand.Text()
	return NewString[:5]
}
