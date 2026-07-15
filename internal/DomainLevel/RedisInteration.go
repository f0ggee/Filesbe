package DomainLevel

import "context"

const (
	RedisPassword = "REDIS_PASSWORD"
	RedisHost     = "REDIS_HOST"
	RedisServer   = "REDIS_SERVER"
)

type DeleterRedis interface {
	DeleteFileInfo(string, context.Context) error
	DeleterFileInfoTest(string, context.Context) error
}

type WriteDataIncomeData struct {
	FileName string
	Info     []byte
	Ctx      context.Context
}

type WritingRedis interface {
	WriteData(WriteDataIncomeData) error
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
