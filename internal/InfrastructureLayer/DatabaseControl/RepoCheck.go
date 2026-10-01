package DatabaseControl

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CheckerDb struct {
	Db *pgxpool.Pool
}

func GetNewCheckerDb(db *pgxpool.Pool) *CheckerDb {
	return &CheckerDb{Db: db}
}

const ErrorUserExist = "an user already exists"

func (db *CheckerDb) CheckerUser(ctx context.Context, email string) error {
	logger := slog.With("CheckUser")

	var existingPerson bool

	err := db.Db.QueryRow(ctx, "SELECT EXISTS (select 1 FROM person WHERE email=$1)", email).Scan(&existingPerson)
	switch {
	case errors.Is(err, context.DeadlineExceeded):
		logger.Error("The context is end", "ERROR", err)
		return errors.New(ErrorTimeEnd)

	case err != nil:
		logger.Info("the strange error", "ERROR", err)
		return err

	}
	if existingPerson {
		return errors.New(ErrorUserExist)
	}

	return nil
}
