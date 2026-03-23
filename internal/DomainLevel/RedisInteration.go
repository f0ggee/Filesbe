package DomainLevel

import "context"

type DeleterRedis interface {
	DeleteFileInfo(string) error
}

type WritingRedis interface {
	WriteData(string, []byte, context.Context) error
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
