package DeletingRedis

import (
	"context"

	"golang.org/x/exp/slog"
)

func (d *DeleterRedis) DeleteFileInfo(nameFileinfo string, ctx context.Context) error {

	err := d.Re.Del(context.Background(), nameFileinfo).Err()
	if err != nil {
		slog.Error("File info's already been deleted", err)
		return err
	}
	return nil
}
