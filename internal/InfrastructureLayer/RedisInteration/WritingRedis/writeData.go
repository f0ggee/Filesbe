package WritingRedis

import (
	"Kaban/internal/Dto"
	"context"
	"log/slog"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func init() {
	err := godotenv.Load("iternal/Service/.env")
	if err != nil {
		slog.Error("cannot load env file", err)
		return

	}
}

func (s *Writing) WriteData(shortName string, InfoAboutFile []byte) error {

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
