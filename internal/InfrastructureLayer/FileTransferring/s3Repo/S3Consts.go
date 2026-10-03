package s3Repo

import "errors"

var (
	ErrorFileSize     = errors.New("size of file is invalid")
	ErrorStrangeError = errors.New("a strange error")
	ErrorClientStop   = errors.New("a client stopped uploading")
	ErrorFileName     = errors.New("the file name isn't defined")
	ErrorNoFile       = errors.New("the needed file was probably deleted or already downloaded")
	ErrorFileDeleted  = errors.New("the file was probably deleted")
)
