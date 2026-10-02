package s3Repo

import "errors"

const ErrorFileSize = "size of file is invalid"

var ErrorClientStop = errors.New("a client stopped uploading")
var ErrorStrangeError = errors.New("a strange error")
var ErrorFileName = "the file name isn't defined"

const ErrorNoFile = "the needed file was probably deleted or already downloaded"
const ErrorFileDeleted = "the file was probably deleted"
