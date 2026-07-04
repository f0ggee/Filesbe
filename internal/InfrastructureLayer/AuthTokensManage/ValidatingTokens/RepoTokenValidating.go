package ValidatingTokens

type Checking struct {
	Key []byte
}

func GetNewChecking(key []byte) *Checking {
	return &Checking{Key: key}
}
