package validator

import (
	"errors"
	"log/slog"
)

func CheckFileSize2(sizeAndName int64) error {

	if sizeAndName > 1*1024*1024*1024 {
		slog.Info("File too big")

		return errors.New("file too big")
	}

	return nil
}
