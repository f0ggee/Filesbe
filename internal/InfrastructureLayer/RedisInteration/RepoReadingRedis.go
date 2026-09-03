package RedisInteration

import (
	"Kaban/internal/DomainLevel"
	"Kaban/internal/Dto"
	"context"
	"errors"
	"log/slog"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	ErrorFindFileInfo = "file's info wasn't found"
)

type RedisReader struct {
	Re *redis.Client
}

func GetNewRedisReader(re *redis.Client) *RedisReader {
	return &RedisReader{Re: re}
}

func (d *RedisReader) GetFileInfo(fileInfoName string, ctx context.Context) ([]byte, error) {

	StructOfFileInfo := Dto.FileInfoLabels{
		InfoAboutFile: nil,
	}

	err := d.Re.HGetAll(ctx, fileInfoName).Scan(&StructOfFileInfo)
	if err != nil {
		slog.Error("Redis GetFileInfo; error happened during getting info about a file", "ERROR", err)
		return nil, errors.New(ErrorFindFileInfo)
	}

	return StructOfFileInfo.InfoAboutFile, nil

}
func (d *RedisReader) GetKey(Ctx context.Context) ([]byte, error) {
	count, sec := 0, 1

	ctx, cancel := context.WithTimeout(Ctx, time.Second*10)
	defer cancel()
	for {
		if count > 20 {
			return nil, errors.New(ErrorReadTimeout)
		}
		err := d.Re.Get(ctx, os.Getenv(DomainLevel.ServerName)).Err()

		if err != nil {
			count, sec = +1, +1
			time.Sleep(time.Duration(sec) * time.Second)
			continue
		}
		var data []byte
		err = d.Re.Get(ctx, os.Getenv(DomainLevel.ServerName)).Scan(&data)
		if err != nil {
			slog.Error("GetKey; error to read data", "ERROR", err)
			return nil, errors.New(ErrorRead)
		}

		err = d.Re.Del(ctx, os.Getenv(DomainLevel.ServerName)).Err()
		if err != nil {
			slog.Error("GetKey; error to delete file info", "ERROR", err)
			return nil, errors.New(ErrorDeleteInfo)
		}
		return data, nil

	}
}
