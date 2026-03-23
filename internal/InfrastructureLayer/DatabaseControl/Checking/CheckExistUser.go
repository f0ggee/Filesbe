package Checking

import (
	"context"
	"errors"

	"golang.org/x/exp/slog"
)

func (db *CheckerDb) CheckerUser(email string, ctx context.Context) error {

	var existingPerson bool

	err := db.Db.QueryRow(ctx, "SELECT EXISTS (select 1 FROM person WHERE email=$1)", email).Scan(&existingPerson)

	switch {

	case errors.Is(err, context.DeadlineExceeded):
		slog.Error("In register,context exceeded", "err - ", err)
		return err

	case err != nil:
		slog.Error("Something wrong in register  check", "err", err)
		return err

	}

	if existingPerson {
		return errors.New("person already exist")
	}

	return nil
}
