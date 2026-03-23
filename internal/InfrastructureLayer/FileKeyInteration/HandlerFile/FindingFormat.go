package HandlerFile

import (
	"mime"
	"path/filepath"
)

func (p ProcessingFile) FindFormatOfFile(s string) string {
	fileExtension := filepath.Ext(s)

	FileExtension := mime.TypeByExtension(fileExtension)
	return FileExtension
}
