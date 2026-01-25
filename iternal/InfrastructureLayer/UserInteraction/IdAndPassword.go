package UserInteraction

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
)

func (db *DB) GetIdPassowrd(email string) (int, string, error) {
	var (
		UnicId   int
		password string
	)

	err := db.Db.QueryRow(context.Background(), `SELECT unic_id ,password FROM person WHERE email=$1`, email).Scan(&UnicId, &password)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		slog.Error("Err don't find anathing in IdAndPassword", "err", err)
		return 0, "", err

	case errors.Is(err, context.DeadlineExceeded):
		slog.Error("Context Dedline Exceed", "Error", err)
		return 0, "", err

	}
	if err != nil {
		slog.Error("Err doesn't equal nil at IdAndPassword ", "Err", err)
		return 0, "", err
	}

	return UnicId, password, nil
}
