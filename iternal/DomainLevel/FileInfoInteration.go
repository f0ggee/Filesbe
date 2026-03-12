package DomainLevel

type FileInfoManipulation interface {
	ConvertToBytesFileInfo(string, []byte) ([]byte, error) //keep it
	GetRealNameFile(string) string                         // Keep it
	ProcessingFileParameters(string) (string, error)       //keep it 0
	GenerateShortFileName() string                         //keep it
	FindFormatOfFile(string) string
	SayHi() string
}
type FileDataManipulation interface {
	EncryptData([]byte, []byte) ([]byte, error)
	DecryptFileInfo([]byte, []byte, []byte) ([]byte, string, error)
}
