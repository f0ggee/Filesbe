package DomainLevel

type FileInfo interface {
	SaveFileInfo(string, string)
	GetRealNameFile(string) string
	ProcessingFileParameters(string) (string, error)
}
