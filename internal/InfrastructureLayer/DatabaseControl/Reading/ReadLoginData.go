package Reading

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
)

func (D Read) LoginData(s string, ctx context.Context) (int, string, error) {
	var (
		id       int
		password string
	)

	err := D.Db.QueryRow(ctx, `SELECT unic_id,password  FROM person WHERE email=$1`, s).Scan(&id, &password)
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
	return id, password, nil
}
