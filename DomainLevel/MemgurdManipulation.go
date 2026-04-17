package DomainLevel

type KeyManage interface {
	SwapingOldKey()
	InstallingNewKey([]byte)
}
