package FileKeyInteration

import (
	"mime"
	"path/filepath"
)

func (*FileInfoController) FindFormatOfFile(FileName string) string {
	fileExtension := filepath.Ext(FileName)

	FileExtension := mime.TypeByExtension(fileExtension)
	return FileExtension

}
