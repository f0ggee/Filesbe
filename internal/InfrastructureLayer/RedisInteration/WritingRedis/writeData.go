package WritingRedis

import (
	"Kaban/internal/Dto"
	"context"
	"log/slog"
)

func (s *Writing) WriteData(shortName string, InfoAboutFile []byte, ctx context.Context) error {

	err := s.Re.HSet(context.Background(), shortName, Dto.FileInfoLabels{
		InfoAboutFile:   InfoAboutFile,
		IsStartDownload: false,
	}).Err()
	if err != nil {
		slog.Error("redis set err", err)
		return err
	}

	return nil

}
