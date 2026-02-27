package RedisInteration

import (
	"Kaban/iternal/Dto"
	"context"
	"errors"
	"log/slog"
	"time"
)

func (d *RedisInterationLayer) GetKey() ([]byte, []byte, []byte, error) {
	count := 0

	for {

		if count > 20 {
			return nil, nil, nil, errors.New("timeout")
		}
		err := d.Re.HGetAll(context.Background(), "server1").Err()

		if err != nil {
			slog.Error("We got the error", "Error", err)
			count++
			time.Sleep(4 * time.Second)
			continue

		}

		zs := Dto.RedisFileStructFromMasterServer{
			AesKey:    nil,
			PlainText: nil,
			Signature: nil,
		}
		err = d.Re.HGetAll(context.Background(), "server1").Scan(&zs)
		if err != nil {
			slog.Error("We got the error when try get the data", "Error", err)
			return nil, nil, nil, errors.New(err.Error())
		}

		// TODO - When i start the project, must remove them
		err = d.Re.Del(context.Background(), "server2").Err()
		if err != nil {
			slog.Error("We got the error", "Error", err)
			return nil, nil, nil, errors.New(err.Error())
		}

		slog.Info("We got the key", "key", true)
		return zs.AesKey, zs.PlainText, zs.Signature, nil

	}
}
