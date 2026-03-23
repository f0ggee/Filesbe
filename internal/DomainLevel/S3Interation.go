package DomainLevel

import "context"

type DeleterS3 interface {
	DeleteFileFromS3(string, context.Context) error
}
