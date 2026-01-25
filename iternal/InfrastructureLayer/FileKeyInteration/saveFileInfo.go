package FileKeyInteration

import (
	"Kaban/iternal/Dto"
	"time"
)

type FileInfoController struct{}

func (*FileInfoController) SaveFileInfo(NameFile, AesKey string) {
	Dto.MapForFile[NameFile] = struct {
		AesKey          string
		TimeSet         time.Time
		IsUsed          bool
		IsStartDownload bool
	}{AesKey: AesKey, TimeSet: time.Now(), IsUsed: false, IsStartDownload: false}

}
