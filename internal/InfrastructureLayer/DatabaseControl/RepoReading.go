package DatabaseControl

import (
	"Kaban/internal/DomainLevel"
	"context"
	"database/sql"
	"errors"
	"log/slog"
)

type Read struct {
	db DatabaseConn
}

func NewRead(db DatabaseConn) *Read {
	return &Read{db: db}
}

var (
	ErrorUserAccount = errors.New("a account doesn't exist")
)

func (D Read) LoginData(ctx context.Context, s string) DomainLevel.OutComingLoginData {
	var (
		id       int32
		password string
	)

	err := D.db.Db.QueryRow(ctx, `SELECT unic_id,password  FROM person WHERE email=$1`, s).Scan(&id, &password)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		slog.Error("LoginData; there isn't an user account", "ERROR", err)
		return DomainLevel.OutComingLoginData{
			Err: ErrorUserAccount,
		}

	case errors.Is(err, context.DeadlineExceeded):
		slog.Error("LoginData; the context is expired", "ERROR", err)
		return DomainLevel.OutComingLoginData{Err: ErrorTimeEnd}

	}
	if err != nil {
		slog.Error("LoginData; an error to get user's data", "ERROR", err)
		return DomainLevel.OutComingLoginData{Err: ErrorStrangeDatabaseError}
	}
	return DomainLevel.OutComingLoginData{
		Id:           id,
		HashPassword: password,
	}
}
