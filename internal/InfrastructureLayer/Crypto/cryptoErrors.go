package Crypto

import "errors"

var (
	ErrorInvalidMode     = errors.New("the crypto mode isn't correct")
	ErrorInvalidData     = errors.New("the data isn't correct ")
	ErrorAuthFail        = errors.New("the encryption key isn't correct")
	ErrorStartEncryption = errors.New("can't start encrypting")
)

var (
	ErrorMakeHash      = errors.New("a hash cannot be created")
	ErrorMakeSign      = errors.New("can't make a sign")
	ErrorStrangeCrypto = errors.New("a strange error happened")
)
var (
	ErrorSignCheck = errors.New("the signature isn't correct")
	ErrorPassword  = errors.New("the password isn't correct")
)
