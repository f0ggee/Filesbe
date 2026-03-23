package Writinig

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Writer struct {
	Db *pgxpool.Pool
}
