package WritingRedis

import (
	"context"
	"log/slog"
)

func (d *Writing) EnableDownloadingParameter(nameOfFileInfo string) error {

	err := d.Re.HSet(context.Background(), nameOfFileInfo, "IsStartDownload", true).Err()
	if err != nil {
		slog.Error("Error set up the labels isStartDownload on true", err)
		return err
	}

	return nil
}
