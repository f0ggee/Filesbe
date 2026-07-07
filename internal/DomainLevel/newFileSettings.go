package DomainLevel

import "time"

const (
	FileMaxSize      = 500000000
	DefaultErrorTime = 12 * time.Hour
)

type SetFileSettings interface {
	FindFormatOfFile(string) string
	FindBestOptions(int64) (int, int)
}
