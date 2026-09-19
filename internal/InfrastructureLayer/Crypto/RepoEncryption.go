package Crypto

import (
	"Kaban/internal/DomainLevel"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"errors"
	"io"
	"log/slog"
)

const ErrorStartEncryption = "can't start encrypting"
const ErrorAuthFail = "the encryption key isn't correct"
const ErrorInvalidData = "the data isn't correct "
const ErrorInvalidMode = "the crypto mode isn't correct"

type RsaEncryption struct {
	privKey           *rsa.PrivateKey
	pubKey            *rsa.PublicKey
	isPrivKeyProvided bool
}

func (r RsaEncryption) Encrypter(bytes []byte) ([]byte, error) {
	if r.isPrivKeyProvided {
		slog.Error("Rsa Encrypter: decrypter mode was provided", "Mode", r.isPrivKeyProvided)
		return nil, errors.New(ErrorInvalidMode)
	}
	return rsa.EncryptOAEP(sha256.New(), rand.Reader, r.pubKey, bytes, nil)
}
func (r RsaEncryption) Decrypt(bytes []byte) ([]byte, error) {
	if !r.isPrivKeyProvided {
		slog.Error("Rsa Encrypter: encrypter mode was provided", "Mode", r.isPrivKeyProvided)
		return nil, errors.New(ErrorInvalidMode)
	}
	return rsa.DecryptOAEP(sha256.New(), rand.Reader, r.privKey, bytes, nil)
}

func (r RsaEncryption) MakeCrypto(bytes []byte, spectra int) (DomainLevel.Crypto, error) {
	if spectra == 1 {
		parsedKey, err := x509.ParsePKCS1PublicKey(bytes)
		if err != nil {
			slog.Error("MakerCrypto: invalid data provided", "ERROR", err)
			return nil, errors.New(ErrorInvalidData)
		}
		return RsaEncryption{
			privKey:           nil,
			pubKey:            parsedKey,
			isPrivKeyProvided: false,
		}, nil
	}
	if spectra == 0 {
		parsedKey, err := x509.ParsePKCS1PrivateKey(bytes)
		if err != nil {
			slog.Error("MakerCrypto: invalid data provided", "ERROR", err)
			return nil, errors.New(ErrorInvalidData)
		}
		return RsaEncryption{
			privKey:           parsedKey,
			pubKey:            nil,
			isPrivKeyProvided: true,
		}, nil
	}
	return nil, errors.New(ErrorInvalidData)
}

type AesEncryption struct {
	aesBlock cipher.AEAD
	nonce    []byte
}

func (a AesEncryption) Encrypter(bytes []byte) ([]byte, error) {
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

func (a AesEncryption) MakeCrypto(bytes []byte) (DomainLevel.Crypto, error) {
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
	nonce := make([]byte, aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return AesEncryption{
		aesBlock: aead,
		nonce:    nonce,
	}, nil
}

type AesCtr struct {
	block   cipher.Stream
	isStart bool
	nonce   []byte
}

func NewAesCtr() *AesCtr {
	return &AesCtr{}
}

func (a AesCtr) Encrypter(bytes []byte) ([]byte, error) {

	if !a.isStart {

		if len(bytes[:aes.BlockSize]) == len(a.nonce) {
		}
		cip := make([]byte, len(bytes)+aes.BlockSize)
		copy(cip, a.nonce)
		a.block.XORKeyStream(cip[aes.BlockSize:], bytes)
		a.isStart = true
		return cip, nil
	}
	a.block.XORKeyStream(bytes, bytes)
	return bytes, nil
}

func (a AesCtr) Decrypt(bytes []byte) ([]byte, error) {
	if !a.isStart && len(bytes) <= aes.BlockSize {
		return nil, errors.New(ErrorInvalidData)
	}
	if !a.isStart {
		cip := make([]byte, len(bytes)-aes.BlockSize)
		a.block.XORKeyStream(cip, bytes[aes.BlockSize:])
		a.isStart = true
		return cip, nil
	}

	a.block.XORKeyStream(bytes, bytes)
	return bytes, nil
}

const ErrorCryptoKey = "the key isn't correct"

func (a AesCtr) MakeCrypto(bytes []byte) (DomainLevel.Crypto, error) {
	if len(bytes) <= aes.BlockSize {
		return nil, errors.New(ErrorCryptoKey)
	}
	block, err := aes.NewCipher(bytes[aes.BlockSize:])
	if err != nil {
		return nil, err
	}
	ctrBlock := cipher.NewCTR(block, bytes[:aes.BlockSize])
	return AesCtr{
		block:   ctrBlock,
		isStart: false,
		nonce:   bytes,
	}, nil
}
