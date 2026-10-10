package RedisInteration

import (
	"Kaban/internal/DomainLevel"
	"context"
	"errors"
	"log/slog"
	"os"
	"time"
)

const (
	ErrorFindFileInfo = "file's info wasn't found"
	ServerName        = "serverName"
)

type RedisRead struct {
}

func NewRedisRead() RedisRead {
	return RedisRead{}
}

func (d RedisRead) GetFileInfo(fileInfoName string, ctx context.Context) ([]byte, error) {

	StructOfFileInfo := DomainLevel.FileInfoLabels{
		InfoAboutFile: nil,
	}

	err := RedisConn.HGetAll(ctx, fileInfoName).Scan(&StructOfFileInfo)
	if err != nil {
		slog.Error("Redis GetFileInfo; error happened during getting info about a file", "ERROR", err)
		return nil, errors.New(ErrorFindFileInfo)
	}

	return StructOfFileInfo.InfoAboutFile, nil

}
func (d RedisRead) GetKey(Ctx context.Context) ([]byte, error) {
	count, sec := 1, 1

	ctx, cancel := context.WithTimeout(Ctx, time.Second*10)
	defer cancel()
	for {
		if count > 20 {
			return nil, errors.New(ErrorReadTimeout)
		}
		err := RedisConn.Get(ctx, os.Getenv(ServerName)).Err()

		if err != nil {
			count, sec = +1, +1
			time.Sleep(time.Duration(sec) * time.Second)
			continue
		}
		var data []byte
		err = RedisConn.Get(ctx, os.Getenv(ServerName)).Scan(&data)
		if err != nil {
			slog.Error("GetPacketData; error to read data", "ERROR", err)
			return nil, errors.New(ErrorRead)
		}

		err = RedisConn.Del(ctx, os.Getenv(ServerName)).Err()
		if err != nil {
			slog.Error("GetPacketData; error to delete file info", "ERROR", err)
			return nil, errors.New(ErrorDeleteInfo)
		}
		return data, nil

	}
}
