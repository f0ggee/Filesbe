package Reading

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Read struct {
	Db *pgxpool.Pool
}

func (D *Read) GetIdPassowrd(email string) (int, string, error) {
	var (
		password string
	)

	err := D.Db.QueryRow(context.Background(), `SELECT unic_id ,password FROM person WHERE email=$1`, email).Scan(&UnitId, &password)

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

	return UnitId, password, nil
}

func (D Read) GetIdPassword(s string) (int, string, error) {
	//T!ODO implement me
	panic("implement me")
}

func (D Read) GetId(s string, ctx context.Context) (int, error) {

	NewCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	var UnitId int

	err := D.Db.QueryRow(NewCtx, `SELECT unic_id  FROM person WHERE email=$1`, s).Scan(&UnitId)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		slog.Error("Err don't find anathing in IdAndPassword", "err", err)
		return 0, err

	case errors.Is(err, context.DeadlineExceeded):
		slog.Error("Context Dedline Exceed", "Error", err)
		return 0, err

	}
	if err != nil {
		slog.Error("Err doesn't equal nil at IdAndPassword ", "Err", err)
		return 0, err
	}

	return UnitId, nil
}

func (D Read) GetPassword(s string, ctx context.Context) (string, error) {

	go func() {

		select {
		case <-ctx.Done():
			return
		}
	}()
	var password string

	err := D.Db.QueryRow(ctx, `SELECT unic_id  FROM person WHERE email=$1`, s).Scan(&password)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		slog.Error("Err don't find anathing in IdAndPassword", "err", err)
		return "", err

	case errors.Is(err, context.DeadlineExceeded):
		slog.Error("Context Dedline Exceed", "Error", err)
		return "", err

	}
	if err != nil {
		slog.Error("Err doesn't equal nil at IdAndPassword ", "Err", err)
		return "", err
	}

	return UnitId, nil

}
