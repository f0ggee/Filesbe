package DatabaseControl

import (
	"Kaban/internal/DomainLevel"
	"context"
	"crypto/rand"
	"database/sql"
	"errors"
	"log/slog"
)

type Read struct {
}

func NewRead() *Read {
	return &Read{}
}

var (
	ErrorUserAccount = errors.New("a account doesn't exist")
)

func (D Read) LoginData(ctx context.Context, s string) DomainLevel.OutComingLoginData {
	var (
		id       int32
		password string
	)

	err := GetConnection().QueryRow(ctx, `SELECT unic_id,password  FROM person WHERE email=$1`, s).Scan(&id, &password)
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

type ReaderTest struct {
}

const testHash = "$2a$10$G2fbJJAq2DRNYSUUjPg9Cu9KBwi/AXraQf6oViu7sE6wuaA4j7Kxa"

func NewReaderTest() *ReaderTest {
	return &ReaderTest{}
}

func (r ReaderTest) LoginData(ctx context.Context, s string) DomainLevel.OutComingLoginData {

	return DomainLevel.OutComingLoginData{
		Id:           123,
		HashPassword: testHash,
		Err:          nil,
	}
}

type ReaderErrorTest struct {
}

func NewReaderErrorTest() *ReaderErrorTest {
	return &ReaderErrorTest{}
}

func (r ReaderErrorTest) LoginData(ctx context.Context, s string) DomainLevel.OutComingLoginData {

	return DomainLevel.OutComingLoginData{
		Err: errors.New(rand.Text()),
	}
}
