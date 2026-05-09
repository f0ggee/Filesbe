package GlobalProces

import (
	"MasterServer_/DomainLevel"
	"fmt"
)

type ProcessController struct {
	Cryptos          DomainLevel.Encryption
	CryptoGen        DomainLevel.CryptoGenerator
	RedisInteracting DomainLevel.RedisUse
	ServerManagement DomainLevel.GettingServersInfo
	CryptoKey        DomainLevel.CryptoKeyManager
	TimeData         *DomainLevel.PreviousSwapTime
}

type ControllingExchange struct {
	E ProcessController
}

func (receiver ControllingExchange) TestTime() {
	fmt.Println("Time", receiver.E.TimeData.GetId())

}

func NewAnotherProcessController(e ProcessController) *ControllingExchange {
	return &ControllingExchange{E: e}
}
