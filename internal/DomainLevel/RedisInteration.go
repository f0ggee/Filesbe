package DomainLevel

type DeleterRedis interface {
	DeleteFileInfo(string) error
}

type WritingRedis interface {
	WriteData(string, []byte) error
	EnableDownloadingParameter(string) error
}

type RedisChecker interface {
	ChekIsStartDownload(string) bool
	CheckFileInfoExists(string) bool
}

type ReadingRedis interface {
	GetKey() ([]byte, []byte, []byte, error)
	GetFileInfo(string) ([]byte, error)
}
