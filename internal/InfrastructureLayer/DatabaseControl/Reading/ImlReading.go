package Reading

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Read struct {
	Db *pgxpool.Pool
}
