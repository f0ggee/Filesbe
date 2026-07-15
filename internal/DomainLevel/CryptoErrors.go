package DomainLevel

const (
	ErrorPassword      = "the password isn't correct"
	ErrorMakeHash      = "a hash cannot be created"
	ErrorStrangeCrypto = "an unexpected error happened during encrypting"
	ErrorSignCheck     = "the signature isn't correct"
	ErrorDecryptKeys   = "data can't be decrypted because of an unexpected error"
	ErrorMakeSign      = "can't make a sign"
	ErrorDataNil       = "the data address is nil"
)
const (
	ErrorDecryptFileInfo = "an error happened during decrypting info"
	ErrorGetFileInfo     = "can't parse file info"
	ErrorInfoOld         = "the file info isn't valid"
	ErrorFileNameEncrypt = "can't encrypt the file because of the file's size"
	ErrorStrangeRead     = "an unexpected error happened during reading a file"
)
