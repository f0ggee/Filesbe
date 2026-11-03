package testFileTransfering

import (
	"Kaban/internal/DomainLevel"
	"errors"
)

type Flow struct {
	U  DomainLevel.MakerUploader
	D  DomainLevel.MakerDownloader
	D_ DomainLevel.MakerDeleter
}

func GetTestFlow() Flow {
	h := &hub{
		data: []byte{},
	}
	return Flow{
		U: testUploader{
			d: h,
		},
		D: testDownload{
			h: h,
		},
		D_: testDelete{
			h: h,
		},
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
