package DeletingRedis

import (
	"Kaban/internal/DomainLevel"
	"context"
	"errors"

	"golang.org/x/exp/slog"
)

func (d *DeleterRedis) DeleteFileInfo(fileInfo string, ctx context.Context) error {

	err := d.Re.Del(ctx, fileInfo).Err()
	if err != nil {
		slog.Error("File info's already been deleted", err)
		return errors.New(DomainLevel.ErrorDeleteInfo)
	}
	return nil
}
