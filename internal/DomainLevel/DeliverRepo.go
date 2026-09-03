package DomainLevel

// The password's errors
const (
	PasswordSizeSmall = "the password is smaller than the default size"
	PasswordSizeBig   = "the password is bigger than the default size"
	PasswordEmpty     = "the password field is empty"
)
const NonIdentifyError = "the strange error"
const ErrorParseInfo = "happened an error during getting data"

const SessionError = "cannot create a session"

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

type LoginApplicationOutComingData struct {
	Jwt string
	Rft string
	Err error
}
type RegisterApplicationOutComingData struct {
	Jwt string
	Rft string
	Err error
}
