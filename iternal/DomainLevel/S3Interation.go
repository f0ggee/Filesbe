package DomainLevel

type S3Interation interface {
	DeleteFileFromS3(string, string) error
}
