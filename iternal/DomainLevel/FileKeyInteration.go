package DomainLevel

type FileInfo interface {
	ConvertToBytesFileInfo(string, []byte) ([]byte, error)
	GetRealNameFile(string) string
	ProcessingFileParameters(string) (string, error)
	GenerateShortFileName() string
	EncryptData([]byte, []byte) ([]byte, error)
	DecryptFileInfo([]byte, []byte, []byte) ([]byte, string, error)
	FindFormatOfFile(string) string
	SayHi() string
}
