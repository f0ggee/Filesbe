package InfrastructureLayer

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	Db *pgxpool.Pool
}

func (db *DB) GeTScrypt(ctx context.Context, UnicId int) (string, error) {
	var scrypt string

	s, err := db.Db.Query(ctx, "SELECT scrypt_salt  FROM person WHERE unic_id = $1", UnicId)
	defer s.Close()
	switch {
	case errors.Is(err, sql.ErrNoRows):
		slog.Error("Doesn't cat select", "err", err)
		return "", err

	case errors.Is(err, context.DeadlineExceeded):
		slog.Error("Context Dedline Exceed", "Error", err)
		return "", err

	}
	if err != nil {
		slog.Error("Err was not equal err while choose", "Err", err)
		return "", err
	}

	for s.Next() {
		err = s.Scan(&scrypt)

	}

	return scrypt, err
}
