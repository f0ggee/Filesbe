package FileControls

import (
	"mime"
	"path/filepath"
)

func FindFormatOfFile(s string) string {
	fileExtension := filepath.Ext(s)

	FileExtension := mime.TypeByExtension(fileExtension)
	return FileExtension
}

func FindBestOptions(size int64) (int, int) {
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
