package FileControls

import (
	"mime"
	"path/filepath"
)

type FileSettings struct{}

func GetNewFileSettings() *FileSettings {
	return &FileSettings{}
}

func (p FileSettings) FindFormatOfFile(s string) string {
	fileExtension := filepath.Ext(s)

	FileExtension := mime.TypeByExtension(fileExtension)
	return FileExtension
}

func (p FileSettings) FindBestOptions(size int64) (int, int) {
	switch {
	case size >= 100*1024*1024:

		fileResult := size / 1000000

		x := 50
		ResultPart := int(fileResult) / x

		NumOfGoroutine := ResultPart + 1
		return x, NumOfGoroutine

	default:
		return 5, 20

	}
}
