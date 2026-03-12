package Dto

type GrpcSendingData struct {
	AesKeyData []byte `json:"aes_key_data"`
	CipherData []byte `json:"cipher_data"`
}
