package HandlingRequests

import "Kaban/internal/DomainLevel"

type HandlerGrpcRequest struct {
	CryptoEncrypt    DomainLevel.Encryption
	CryptoDecrypt    DomainLevel.Decryption
	CryptoValidate   DomainLevel.CryptoValidating
	ValidationPacket DomainLevel.PacketChecker
}
