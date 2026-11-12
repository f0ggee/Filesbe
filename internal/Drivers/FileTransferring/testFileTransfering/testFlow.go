package testFileTransfering

import (
	"Kaban/internal/DomainLevel"
	"errors"
)

type Flow struct {
	U  DomainLevel.MakerUploader
	D  DomainLevel.MakerDownloader
	D_ DomainLevel.MakerDeleter
	h  *[]byte
}

func (f *Flow) Check() []byte {
	return *f.h
}

func (f *Flow) Reset() {
	*f.h = nil
}

func GetTestFlow() *Flow {
	var testHub []byte
	return &Flow{
		U: testUploader{
			&testHub,
		},
		D: testDownload{
			&testHub,
		},
		D_: testDelete{
			h: &testHub,
		},
		h: &testHub,
	}
}

func GetTestErrorFlow() Flow {
	return Flow{
		U: ErrorUpload{err: errors.New("Error")},
		D: DownloadError{err: errors.New("Errror")},
		D_: DeleteError{
			err: errors.New("error"),
		},
	}
}
