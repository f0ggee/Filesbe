package Application

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/Dto"
	"context"
	"crypto/rand"
	"crypto/x509"
	"log/slog"
	"os"
	"time"

	"github.com/awnumar/memguard"
	"golang.org/x/sync/errgroup"
)

type NewSwapKeyFirst struct {
	GetCrypto
	GetControlKeys
	Parser
}

func GetNewNewSwapKeyFirst(crypto GetCrypto, controlKeys GetControlKeys, parser Parser) *NewSwapKeyFirst {
	return &NewSwapKeyFirst{GetCrypto: crypto, GetControlKeys: controlKeys, Parser: parser}
}

func (sa *NewSwapKeyFirst) GetSwapKeyFirst() time.Duration {

	return MakerRequests(sa, convertedDataGrpcDataLooks)
}

func (sa *NewSwapKeyFirst) MakerRequests(convertedDataGrpcDataLooks []byte) time.Duration {

}
