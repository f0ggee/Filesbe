package DomainLevel

type RedisInteration interface {
	WriteData(string, []byte) error
	ChekIsStartDownload(string) bool
	SetIstartDonwload(string) error
	GetFileInfo(string) ([]byte, error)
	DeleteFileInfo(string) error
}
