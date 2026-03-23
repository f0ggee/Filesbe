package Checking

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type CheckerDb struct {
	Db *pgxpool.Pool
}
