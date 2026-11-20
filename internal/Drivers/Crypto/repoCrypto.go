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
	"io"
	"log/slog"
)

type RsaEncryption struct {
	privKey *rsa.PrivateKey
	pubKey  *rsa.PublicKey
	kind    []byte
}

func (r *RsaEncryption) GetRequiredOverheadSize() int {
	return 0
}

func NewRsaEncryption() *RsaEncryption {
	return &RsaEncryption{}
}

func (r *RsaEncryption) Encrypt(bytes []byte) ([]byte, error) {
	if r.pubKey == nil {
		return nil, ErrorInvalidMode
	}
	return rsa.EncryptOAEP(sha256.New(), rand.Reader, r.pubKey, bytes, nil)
}
func (r *RsaEncryption) Decrypt(bytes []byte) ([]byte, error) {
	if r.privKey == nil {
		return nil, ErrorInvalidMode
	}

	return rsa.DecryptOAEP(sha256.New(), rand.Reader, r.privKey, bytes, nil)
}

func (r *RsaEncryption) MakeCrypto(key []byte, spectr []byte) (DomainLevel.Crypto, error) {
	var err error
	parsedKey, err := x509.ParsePKCS1PrivateKey(key)
	if err == nil {
		return &RsaEncryption{
			privKey: parsedKey,
			pubKey:  &parsedKey.PublicKey,
		}, nil
	}

	publicKey, err := x509.ParsePKCS1PublicKey(key)
	if err == nil {
		return &RsaEncryption{
			privKey: nil,
			pubKey:  publicKey,
		}, nil
	}

	slog.Error("RsaCrypotMaker: error to determine a key", "ERROR", err)
	return nil, ErrorInvalidData
}

type AesEncryption struct {
	aesBlock       cipher.AEAD
	nonce          []byte
	isStartEncrypt bool
	isNonceGet     bool
}

func NewAesGcmEncryption() AesEncryption {
	return AesEncryption{}
}

func (a AesEncryption) GetRequiredOverheadSize() int {
	return 0
}

var ErrorCryptoInvalidData = errors.New("the crypto input data isn't correct")

var ErrorCrypto = errors.New("the unexpected error")

func (a AesEncryption) Encrypt(data []byte) ([]byte, error) {

	nonce := make([]byte, a.aesBlock.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, ErrorCrypto
	}

	return a.aesBlock.Seal(nonce, nonce, data, nil), nil
}

func (a AesEncryption) Decrypt(cipherText []byte) ([]byte, error) {
	if len(cipherText) <= a.aesBlock.NonceSize()+a.aesBlock.Overhead() {
		return nil, ErrorCryptoInvalidData
	}
	nonce := cipherText[:a.aesBlock.NonceSize()]
	data := cipherText[a.aesBlock.NonceSize():]
	return a.aesBlock.Open(nil, nonce, data, nil)
}

func (a AesEncryption) MakeCrypto(bytes []byte, i []byte) (DomainLevel.Crypto, error) {
	block, err := aes.NewCipher(bytes)
	if err != nil {
		slog.Error("MakeCrypto AES: error can't made a block", "ERROR", err)
		return nil, ErrorStartEncryption
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		slog.Error("MakeCrypto AES: error can't initializer a GCM mode", "ERROR", err)
		return nil, ErrorStartEncryption
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

func (a *AesCtrEncryption) GetRequiredOverheadSize() int {

	return aes.BlockSize
}

func NewAesCtr() *AesCtrEncryption {
	return &AesCtrEncryption{}
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
		return nil, ErrorCryptoInvalidData
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	ctrBlock := cipher.NewCTR(block, i)
	a.nonce = i
	a.block = ctrBlock
	a.isStart = false
	return a, nil
}
