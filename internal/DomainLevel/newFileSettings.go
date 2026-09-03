package DomainLevel

type SetFileSettings interface {
	FindFormatOfFile(string) string
	FindBestOptions(int64) (int, int)
}
