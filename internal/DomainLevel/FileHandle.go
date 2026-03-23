package DomainLevel

type HandleFileInfo interface {
	GetRealNameFile(string) string // Keep it
	//ProcessingFileParameters(string) (string, error)       //keep it 0
	SayHi() string
}

type HandleFile interface {
	FindFormatOfFile(string) string
	FindBesOptions(int64) (int, int)
}
