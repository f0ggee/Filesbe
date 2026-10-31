package DatabaseControl

import (
	"Kaban/internal/DomainLevel"
	"context"
	"crypto/rand"
	"errors"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
)

type Writer struct {
}

func NewWriter() *Writer {
	return &Writer{}
}

func (d *Writer) CreateUser(data DomainLevel.CreateUserIncomingData) (int32, error) {
	var UnitId int32
	tx, err := GetConnection().Begin(data.Ctx)
	if err != nil {
		slog.Error("CreateUser; error to start a transaction", "ERROR", err)
		return 0, ErrorStrangeDatabaseError
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil {
			slog.Error("CreateUser; the transaction wasn't finished", "ERROR", err)
		}
	}(tx, data.Ctx)

	err = tx.QueryRow(data.Ctx, "INSERT INTO person(name,email,password,created_at,scrypt_salt) VALUES ($1,$2,$3,$4,$5) RETURNING unic_id", data.Name, data.Email, data.HashPassword, time.Now(), rand.Text()).Scan(&UnitId)

	switch {
	case errors.Is(err, context.DeadlineExceeded):
		slog.Error("CreateUser; the context is expired", "ERROR", err)
		return 0, ErrorTimeEnd

	case err != nil:
		slog.Error("CreateUser; a strange error", "ERROR", err)
		return 0, ErrorStrangeDatabaseError
	}

	if err = tx.Commit(data.Ctx); err != nil {
		slog.Error("CreateUser:error can't commit", "ERROR", err)
		return 0, err
	}
	return UnitId, nil

}
