package RedisInteration

import "errors"

var (
	ErrorDeleteInfo   = errors.New("an unexpected error happened during deleting a file")
	ErrorWrite        = errors.New("an unexpected error happened during file writing")
	ErrorRead         = errors.New("an unexpected error happened during file reading")
	ErrorReadTimeout  = errors.New("the time of reading is finished")
	ErrorFindFileInfo = errors.New("file's info wasn't found")
)
