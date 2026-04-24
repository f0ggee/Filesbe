package WritingRedis

import (
	"context"
	"log/slog"
)

func (d *Writing) EnableDownloadingParameter(nameOfFileInfo string, ctx context.Context) error {

	err := d.Re.HSet(context.Background(), nameOfFileInfo, "IsStartDownload", true).Err()
	if err != nil {
		slog.Error("Error set up the labels isStartDownload on true", "error", err.Error())
		return err
	}

	return nil
}
