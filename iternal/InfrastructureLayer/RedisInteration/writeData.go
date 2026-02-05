package RedisInteration

import (
	"Kaban/iternal/Dto"
	"context"
	"log/slog"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load("iternal/Service/.env")
	if err != nil {
		slog.Error("cannot load env file", err)
		return

	}
}

type RedisInterationLayer struct{}

func (*RedisInterationLayer) WriteData(shortName string, InfoAboutFile []byte) error {

	redisConnect := ConnectToRedis()
	defer redisConnect.Close()
	err := redisConnect.HSet(context.Background(), shortName, Dto.FileInfo{
		InfoAboutFile:   InfoAboutFile,
		IsStartDownload: false,
	}).Err()
	if err != nil {
		slog.Error("redis set err", err)
		return err
	}

	return nil

}
