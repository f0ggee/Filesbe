package Crypto

import (
	"Kaban/internal/DomainLevel"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"errors"
	"log/slog"
)

const ErrorStartEncryption = "can't start encrypting"
const ErrorAuthFail = "the encryption key1 isn't correct"
const ErrorInvalidData = "the data isn't correct "
const ErrorInvalidMode = "the crypto mode isn't correct"

type RsaEncryption struct {
	privKey           *rsa.PrivateKey
	pubKey            *rsa.PublicKey
	isPrivKeyProvided bool
	kind              []byte
}

func NewRsaEncryption() RsaEncryption {
	return RsaEncryption{}
}

func (r *RsaEncryption) GetRequiredOverheadSize() int {
	return 0
}

func (r *RsaEncryption) Encrypt(bytes []byte) ([]byte, error) {
	if r.pubKey == nil {
		slog.Error("Rsa Encrypt: encrypter mode was provided", "Mode", r.isPrivKeyProvided)
		return nil, errors.New(ErrorInvalidMode)
	}
	return rsa.EncryptOAEP(sha256.New(), rand.Reader, r.pubKey, bytes, nil)
}
func (r *RsaEncryption) Decrypt(bytes []byte) ([]byte, error) {
	if r.privKey == nil {
		slog.Error("Rsa Encrypt: decrypter mode was provided", "Mode", r.isPrivKeyProvided)
		return nil, errors.New(ErrorInvalidMode)
	}

	return rsa.DecryptOAEP(sha256.New(), rand.Reader, r.privKey, bytes, nil)
}

func (r *RsaEncryption) MakeCrypto(key []byte, spectr []byte) (DomainLevel.Crypto, error) {
	var err error
	parsedKey, err := x509.ParsePKCS1PrivateKey(key)
	if err == nil {
		return &RsaEncryption{
			privKey:           parsedKey,
			pubKey:            nil,
			isPrivKeyProvided: true,
		}, nil
	}

	publicKey, err := x509.ParsePKCS1PublicKey(key)
	if err == nil {
		return &RsaEncryption{
			privKey:           nil,
			pubKey:            publicKey,
			isPrivKeyProvided: false,
		}, nil
	}

	slog.Error("RsaCrypotMaker: error to determine a key", "ERROR", err)
	return nil, errors.New(ErrorInvalidData)
}

type AesEncryption struct {
	aesBlock cipher.AEAD
	nonce    []byte
}

func NewAesGcmEncryption() AesEncryption {
	return AesEncryption{}
}

func (a AesEncryption) GetRequiredOverheadSize() int {
	return 12
}

const ErrorCryptoInvalidData = "the crypto input data isn't correct"

func (a AesEncryption) Encrypt(bytes []byte) ([]byte, error) {

	return a.aesBlock.Seal(a.nonce, a.nonce, bytes, nil), nil
}

func (a AesEncryption) Decrypt(bytes []byte) ([]byte, error) {

	data, err := a.aesBlock.Open(nil, bytes[:a.aesBlock.NonceSize()], bytes[a.aesBlock.NonceSize():], nil)
	if err != nil {
		slog.Error("Error AES:decrypt data", "ERROR", err)
		return nil, errors.New(ErrorAuthFail)
	}
	return data, nil
}

func (a AesEncryption) MakeCrypto(bytes []byte, i []byte) (DomainLevel.Crypto, error) {
	if len(i) != 12 {
		return nil, errors.New(ErrorCryptoInvalidData)
	}
	block, err := aes.NewCipher(bytes)
	if err != nil {
		slog.Error("MakeCrypto AES: error can't made a block", "ERROR", err)
		return nil, errors.New(ErrorStartEncryption)
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		slog.Error("MakeCrypto AES: error can't initializer a GCM mode", "ERROR", err)
		return nil, errors.New(ErrorStartEncryption)
	}
	return AesEncryption{
		aesBlock: aead,
		nonce:    i,
	}, nil
}

type AesCtrEncryption struct {
	block            cipher.Stream
	isStart          bool
	isDecryptedStart bool
	nonce            []byte
}

func (a AesCtrEncryption) GetRequiredOverheadSize() int {
	return aes.BlockSize
}

func NewAesCtr() AesCtrEncryption {
	return AesCtrEncryption{}
}

func (a *AesCtrEncryption) Encrypt(bytes []byte) ([]byte, error) {

	if !a.isStart {
		cip := make([]byte, len(bytes)+aes.BlockSize)
		copy(cip, a.nonce)
		a.block.XORKeyStream(cip[aes.BlockSize:], bytes)
		a.isStart = true
		return cip, nil
	}
	a.block.XORKeyStream(bytes, bytes)
	return bytes, nil
}

func (a *AesCtrEncryption) Decrypt(cipherText []byte) ([]byte, error) {
	if !a.isDecryptedStart && len(cipherText) > 16 && bytes.Equal(cipherText[:aes.BlockSize], a.nonce) {
		ciphertext := cipherText[aes.BlockSize:]
		a.block.XORKeyStream(ciphertext, ciphertext)
		a.isDecryptedStart = true
		return ciphertext, nil
	}
	a.block.XORKeyStream(cipherText, cipherText)
	return cipherText, nil
}

func (a *AesCtrEncryption) MakeCrypto(key []byte, i []byte) (DomainLevel.Crypto, error) {
	if key == nil || i == nil {
		return nil, errors.New(ErrorCryptoInvalidData)
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	ctrBlock := cipher.NewCTR(block, i)
	return &AesCtrEncryption{
		block:   ctrBlock,
		isStart: false,
		nonce:   i,
	}, nil
}
