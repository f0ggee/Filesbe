package UserInteraction

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
)

func (d *DB) CreateUser(Name, Email, HashPassword, ScryptKey string) (int, error) {

	var unicId int

	tx, err := d.Db.Begin(context.Background())
	if err != nil {
		slog.Error("Error in start transaction", "Err", err)
		return 0, err
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil {
			slog.Error("Error is doing the rollback", "Err", err)

		}
	}(tx, context.Background())

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
