package DatabaseControl

import (
	"Kaban/internal/DomainLevel"
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Read struct {
	Db *pgxpool.Pool
}

const ErrorUserAccount = "a account doesn't exist"
const ErrorTimeEnd = "the creating time is expired"
const ErrorStrangeDatabaseError = "an unexpected error happened during getting user's data"

func GetNewRead(db *pgxpool.Pool) *Read {
	return &Read{Db: db}
}

func (D Read) LoginData(ctx context.Context, s string) DomainLevel.OutComingLoginData {
	var (
		id       int32
		password string
	)

	err := D.Db.QueryRow(ctx, `SELECT unic_id,password  FROM person WHERE email=$1`, s).Scan(&id, &password)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		slog.Error("LoginData; there isn't an user account", "ERROR", err)
		return DomainLevel.OutComingLoginData{
			Err: errors.New(ErrorUserAccount),
		}

	case errors.Is(err, context.DeadlineExceeded):
		slog.Error("LoginData; the context is expired", "ERROR", err)
		return DomainLevel.OutComingLoginData{Err: errors.New(ErrorTimeEnd)}

	}
	if err != nil {
		slog.Error("LoginData; an error to get user's data", "ERROR", err)
		return DomainLevel.OutComingLoginData{Err: errors.New(ErrorStrangeDatabaseError)}
	}
	return DomainLevel.OutComingLoginData{
		Id:           id,
		HashPassword: password,
	}
}
