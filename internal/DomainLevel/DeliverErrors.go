package DomainLevel

// The password's errors
const (
	PasswordSizeSmall = "the password is smaller than the default size"
	PasswordSizeBig   = "the password is bigger than the default size"
	PasswordEmpty     = "the password field is empty"
)
const NonIdentifyError = "the strange error"
const ErrorParseInfo = "happened an error during getting data"

// The email's errors
const (
	NotCorrectEmail = "the email isn't correct"
	EmailEmpty      = "the email's filed is empty"
	EmailMaxSize    = "the email's size is bigger than the default size"
)
const (
	NameMaxSize = "the name size is bigger than the default size"
	NameMinSize = "the name size is smaller then the default size"
)

// Operations' errors
const (
	Break    = "break"
	NotStart = "not_start"
)
const (
	ErrorUserNotAuthed = "user's auth is expired"
	ErrorUserToken     = "user's auth isn't valid"
)

const (
	MethodNotAllowed = "the method isn't allowed"
)
const (
	ErrorFileNameEmpty   = "the file's filed is empty"
	ErrorCantGetFileName = "the file's name isn't set"
)
