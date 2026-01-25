package FileKeyInteration

import (
	Dto2 "Kaban/iternal/Dto"
	"errors"
	"log/slog"
)

func (*FileInfoController) ProcessingFileParameters(RealNameFile string) (string, error) {
	sa, ok := Dto2.MapForFile[RealNameFile]

	if !ok {
		slog.Info("Name file ", RealNameFile)
		slog.Error("File don't find ")
		return "", errors.New("File don't find ")
	}

	if sa.IsUsed {
		slog.Error("File's already used")
		return "", errors.New("file's already used")
	}

	sa.IsStartDownload = true
	return sa.AesKey, nil
}
