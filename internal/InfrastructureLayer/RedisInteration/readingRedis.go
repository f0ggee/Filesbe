package RedisInteration

import (
	"Kaban/internal/DomainLevel"
	"context"
	"crypto/rand"
	"encoding/json"
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
	sec := 1

	ctx, cancel := context.WithTimeout(Ctx, time.Second*10)
	defer cancel()
	for count := range 20 {
		err := RedisConn.Get(ctx, os.Getenv(ServerName)).Err()

		if err != nil {
			sec++
			slog.Group("Redis GetFileInfo:Retry again", slog.Int("Counts:", count), slog.Int("Sec", sec))
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
	return nil, errors.New(ErrorReadTimeout)
}

type RedisReadTest struct {
}

func NewRedisReadTest() *RedisReadTest {
	return &RedisReadTest{}
}

func (r RedisReadTest) GetKey(ctx context.Context) ([]byte, error) {

	d, _ := json.Marshal([]byte(rand.Text()))
	return d, nil
}

func (r RedisReadTest) GetFileInfo(s string, ctx context.Context) ([]byte, error) {
	d, _ := json.Marshal([]byte(rand.Text()))

	return d, nil
}

type RedisTestIO struct {
	d []byte
}

func NewRedisTestIO() *RedisTestIO {
	return &RedisTestIO{}
}

func (r *RedisTestIO) WriteData(data DomainLevel.WriteDataIncomeData) error {
	r.d = data.Info
	return nil
}

func (r *RedisTestIO) EnableDownloadingParameter(s string, ctx context.Context) error {
	return nil
}

func (r *RedisTestIO) GetKey(ctx context.Context) ([]byte, error) {
	return nil, nil
}

func (r *RedisTestIO) GetFileInfo(s string, ctx context.Context) ([]byte, error) {
	return r.d, nil
}
