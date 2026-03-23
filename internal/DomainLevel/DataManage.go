package DomainLevel

type DataConvert interface {
	JsonConverter(any) ([]byte, error)
}
