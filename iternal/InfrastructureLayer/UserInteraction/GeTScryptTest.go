package UserInteraction

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

func (d *DbForTests) GeTScryptTest(unicId int) (string, error) {

	se := rand.Text()

	sa, _ := hex.DecodeString(se)
	hash := sha256.Sum256(sa)
	return hex.EncodeToString(hash[:]), nil
}
