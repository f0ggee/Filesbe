package Writinig

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
)

func (d *Writer) CreateUser(Name string, Email string, HashPassword string, ctx context.Context) (int, error) {

	var UnitId int

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

	err = tx.QueryRow(ctx, "INSERT INTO person(name,email,password,created_at) VALUES ($1,$2,$3,$4,$5) RETURNING unic_id", Name, Email, HashPassword, time.Now()).Scan(&UnitId)

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

	return UnitId, nil

}
