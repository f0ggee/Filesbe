package DomainLevel

const (
	ErrorPassword        = "the password isn't correct"
	ErrorMakeHash        = "a hash cannot be created"
	ErrorStrangeCrypto   = "an unexpected error happened during encrypting"
	ErrorSignCheck       = "the signature isn't correct"
	ErrorDecryptKeys     = "data can't be decrypted because of an unexpected error"
	ErrorDecryptFileInfo = "an error happened during decrypting info"
	ErrorInfoOld         = "the file info isn't valid"
	ErrorGetFileInfo     = "can't parse file info"
	ErrorFileNameEncrypt = "can't encrypt the file because of the file's size"
	ErrorMakeSign        = "can't make a sign"
	ErrorDataNil         = "the data address is nil"
)
