package Validator

import (
	"Kaban/internal/DomainLevel"
	"context"
	"errors"

	"golang.org/x/exp/slog"
)

func (db *CheckerDb) CheckerUser(email string, ctx context.Context) error {
	var existingPerson bool
	err := db.Db.QueryRow(ctx, "SELECT EXISTS (select 1 FROM person WHERE email=$1)", email).Scan(&existingPerson)
	switch {
	case errors.Is(err, context.DeadlineExceeded):
		slog.Error("CheckUser; the context is expired", "ERROR", err)
		return errors.New(DomainLevel.ErrorTimeEnd)

	case err != nil:
		slog.Error("CheckUser; a strange error", "ERROR", err)
		return errors.New(DomainLevel.ErrorStrangeDatabaseError)

	}
	if existingPerson {
		return errors.New(DomainLevel.ErrorUserExist)
	}
	return nil
}
