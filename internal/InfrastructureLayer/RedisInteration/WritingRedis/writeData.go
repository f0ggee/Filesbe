package WritingRedis

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/Dto"
	"context"
	"errors"
	"log/slog"
)

func (s *Writing) WriteData(shortName string, InfoAboutFile []byte, ctx context.Context) error {

	err := s.Re.HSet(ctx, shortName, Dto.FileInfoLabels{
		InfoAboutFile:   InfoAboutFile,
		IsStartDownload: false,
	}).Err()
	if err != nil {
		slog.Error("Redis WriteData; error to write data", "ERROR", err)
		return errors.New(DomainLevel.ErrorWrite)
	}
	return nil

}
