package testFileTransfering

import (
	"Kaban/internal/DomainLevel"
	"errors"
)

type Flow struct {
	U  DomainLevel.MakerUploader
	D  DomainLevel.MakerDownloader
	D_ DomainLevel.MakerDeleter
	h  hub
}

func (f *Flow) Reset() {

	h := &hub{
		data: []byte{},
	}
	f.U = testUploader{
		d: h,
	}
	f.D = testDownload{
		h: h,
	}
	f.D_ = testDelete{
		h: h,
	}
}

func GetTestFlow() *Flow {
	h := &hub{
		data: []byte{},
	}
	return &Flow{
		U: testUploader{
			d: h,
		},
		D: testDownload{
			h: h,
		},
		D_: testDelete{
			h: h,
		},
		h: *h,
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
