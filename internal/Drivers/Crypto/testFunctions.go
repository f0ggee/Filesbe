package Crypto

import (
	"Kaban/internal/DomainLevel"
	"crypto/aes"
	"errors"
	mrand "math/rand/v2"
)

type RsaError struct {
}

func NewRsaError() *RsaError {
	return &RsaError{}
}

func (r RsaError) Encrypt(i []byte) ([]byte, error) {
	return nil, encryptErrors[mrand.IntN(len(encryptErrors))]
}

func (r RsaError) Decrypt(i []byte) ([]byte, error) {
	return nil, encryptErrors[mrand.IntN(len(encryptErrors))]
}

func (r RsaError) MakeCrypto(key []byte, ACD []byte) (DomainLevel.Crypto, error) {
	return r, nil
}

func (r RsaError) GetRequiredOverheadSize() int {
	return 0
}

type CheckingTest struct {
}

func NewCheckingTest() *CheckingTest {
	return &CheckingTest{}
}

func (c CheckingTest) PasswordVerify(bytes []byte, bytes2 []byte) error {
	return nil
}

type CheckingTestError struct {
}

func NewCheckingTestError() *CheckingTestError {
	return &CheckingTestError{}
}

func (c CheckingTestError) PasswordVerify(bytes []byte, bytes2 []byte) error {

	return ErrorPassword
}

type NewAesCtrError struct {
}

var encryptErrors = []error{
	ErrorCrypto,
	ErrorCryptoInvalidData,
	ErrorInvalidMode,
	ErrorInvalidData,
	ErrorAuthFail,
	ErrorStartEncryption,
	ErrorStrangeCrypto,
}

func (n NewAesCtrError) Encrypt(i []byte) ([]byte, error) {
	return nil, encryptErrors[mrand.IntN(len(encryptErrors))]
}

func (n NewAesCtrError) Decrypt(i []byte) ([]byte, error) {
	return nil, encryptErrors[mrand.IntN(len(encryptErrors))]
}

func (n NewAesCtrError) MakeCrypto(key []byte, ACD []byte) (DomainLevel.Crypto, error) {
	return n, nil
}

func (n NewAesCtrError) GetRequiredOverheadSize() int {
	return aes.BlockSize
}

type AesGcmError struct{}

func NewAesGcmError() *AesGcmError {
	return &AesGcmError{}
}

func (a AesGcmError) Encrypt(bytes []byte) ([]byte, error) {
	return nil, errors.New("error")
}

func (a AesGcmError) Decrypt(bytes []byte) ([]byte, error) {
	return nil, nil
}

func (a AesGcmError) MakeCrypto(key []byte, ACD []byte) (DomainLevel.Crypto, error) {

	return a, nil
}

func (a AesGcmError) GetRequiredOverheadSize() int {
	return 16
}
