package DatabaseControl

import (
	"context"
	"errors"
	"log/slog"
)

type CheckerDb struct {
	DatabaseConn
}

func NewCheckerDb(databaseConn DatabaseConn) *CheckerDb {
	return &CheckerDb{DatabaseConn: databaseConn}
}

var ErrorUserExist = errors.New("the user already exists")

func (db *CheckerDb) CheckerUser(ctx context.Context, email string) error {
	logger := slog.With("CheckUser")
	var existingPerson bool
	err := db.Db.QueryRow(ctx, "SELECT EXISTS (select 1 FROM person WHERE email=$1)", email).Scan(&existingPerson)
	switch {
	case errors.Is(err, context.DeadlineExceeded):
		logger.Error("The context is end", "ERROR", err)
		return ErrorTimeEnd

	case err != nil:
		logger.Info("the strange error", "ERROR", err)
		return err

	}
	if existingPerson {
		return ErrorUserExist
	}

	return nil
}
