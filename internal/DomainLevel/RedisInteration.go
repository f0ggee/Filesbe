package DomainLevel

import "context"

type DeleterRedis interface {
	DeleteFileInfo(string, context.Context) error
	DeleterFileInfoTest(string, context.Context) error
}

type WritingRedis interface {
	WriteData(string, []byte, context.Context) error
	EnableDownloadingParameter(string, context.Context) error
}

type RedisChecker interface {
	ChekIsStartDownload(string, context.Context) bool
	ChekIsStartDownloadTest(string, context.Context) bool
	CheckFileInfoExists(string, context.Context) bool
}

type ReadingRedis interface {
	GetKey(context.Context) ([]byte, error)
	GetFileInfo(string, context.Context) ([]byte, error)
}
