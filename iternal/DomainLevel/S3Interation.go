package DomainLevel

type S3Handle interface {
	DeleteFileFromS3(string, string) error
}
