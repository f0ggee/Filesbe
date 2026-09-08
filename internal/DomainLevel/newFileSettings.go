package DomainLevel

import (
	"mime"
	"path/filepath"
)

type FileSettings struct {
	Size int64
	Name string
}

func GetNewFileSettings(size int64, name string) *FileSettings {
	return &FileSettings{Size: size, Name: name}
}

func (p FileSettings) FindFormatOfFile() string {
	fileExtension := filepath.Ext(p.Name)

	FileExtension := mime.TypeByExtension(fileExtension)
	return FileExtension
}

func (p FileSettings) FindBestOptions() (int, int) {
	switch {
	case p.Size >= 100*1024*1024:

		fileResult := p.Size / 1000000

		x := 50
		ResultPart := int(fileResult) / x

		NumOfGoroutine := ResultPart + 1
		return x, NumOfGoroutine

	default:
		return 5, 20

	}
}
