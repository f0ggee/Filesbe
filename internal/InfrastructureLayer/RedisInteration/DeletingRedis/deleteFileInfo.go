package DeletingRedis

import (
	"context"

	"golang.org/x/exp/slog"
)

func (d *DeleterRedis) DeleteFileInfo(nameFileinfo string) error {

	err := d.Re.Del(context.Background(), nameFileinfo).Err()
	if err != nil {
		slog.Error("Error in func deleteFileInfo in Redis", err)
		return err
	}
	return nil
}
