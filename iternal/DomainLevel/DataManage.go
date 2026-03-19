package DomainLevel

type DataConvert interface {
	JsonConverter(any) ([]byte, error)
	ConvertFileInfo(string, []byte) ([]byte, error)
}
