package DomainLevel

type DeleterS3 interface {
	DeleteFileFromS3(string, string) error
}
