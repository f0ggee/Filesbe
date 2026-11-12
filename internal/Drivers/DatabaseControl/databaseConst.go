package DatabaseControl

import "errors"

var (
	ErrorTimeEnd              = errors.New("the creating time is expired")
	ErrorStrangeDatabaseError = errors.New("an unexpected error happened during getting user's data")
)
