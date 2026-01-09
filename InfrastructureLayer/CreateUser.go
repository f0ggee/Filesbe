package InfrastructureLayer

import (
	"context"
	"errors"
	"log/slog"
	"time"
)

func (d *DB) CreateUser(Name, Email, HashPassword, ScryptKey string) (int, error) {

	var unicId int

	tx, err := d.Db.Begin(context.Background())
	if err != nil {
		slog.Error("Error in start transaction", "Err", err)
		return 0, err
	}
	defer tx.Rollback(context.Background())

	err = tx.QueryRow(context.Background(), "INSERT INTO person(name,email,password,scrypt_salt,created_at) VALUES ($1,$2,$3,$4,$5) RETURNING unic_id", Name, Email, HashPassword, ScryptKey, time.Now()).Scan(&unicId)

	switch {
	case errors.Is(err, context.DeadlineExceeded):
		slog.Error("Context exceeded", "Err", err)
		return 0, err

	case err != nil:
		slog.Error("Error something wrong in create user", "Errr", err)
		return 0, err
	}

	if err = tx.Commit(context.Background()); err != nil {
		slog.Error("Error cant commit", "Err", err)
		return 0, err
	}

	return unicId, nil

}
