package Connect_to_BD

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var dbIp = os.Getenv("POSTGRESQL_HOST")
var dbPort = os.Getenv("POSTGRESQL_PORT")
var dbUser = os.Getenv("POSTGRESQL_USER")
var dbPassword = os.Getenv("POSTGRESQL_PASSWORD")
var dbDbname = os.Getenv("POSTGRESQL_DBNAME")

func init() {

	err := godotenv.Load("iternal/Service/.env")
	if err != nil {
		slog.Info("Error loading the file ")
		return
	}

	dbIp = os.Getenv("POSTGRESQL_HOST")
	dbPort = os.Getenv("POSTGRESQL_PORT")
	dbUser = os.Getenv("POSTGRESQL_USER")
	dbPassword = os.Getenv("POSTGRESQL_PASSWORD")
	dbDbname = os.Getenv("POSTGRESQL_DBNAME")

}
func config() *pgxpool.Config {

	const Maxconns = int32(5)
	const Mincons = int32(2)
	const Lifetime = time.Hour
	const IdelTime = time.Minute * 20
	const Health = time.Minute

	_ = godotenv.Load()

	connstr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s", dbIp, dbPort, dbUser, dbPassword, dbDbname)

	dbconfige, err := pgxpool.ParseConfig(connstr)
	if err != nil {
		slog.Error("Error loading database connection", err)
		return nil

	}
	dbconfige.MaxConns = Maxconns
	dbconfige.MinConns = Mincons
	dbconfige.MaxConnLifetime = Lifetime
	dbconfige.MaxConnIdleTime = IdelTime
	dbconfige.HealthCheckPeriod = Health

	dbconfige.BeforeAcquire = func(ctx context.Context, conn *pgx.Conn) bool {

		slog.Info("Settings are set")
		return true
	}
	dbconfige.AfterRelease = func(conn *pgx.Conn) bool {
		slog.Info("After coonect all okay")
		return true

	}

	dbconfige.BeforeClose = func(conn *pgx.Conn) {
		slog.Info("Connection are close ")

	}

	return dbconfige
}

func Connect() (*pgxpool.Pool, error) {

	connPool, err := pgxpool.NewWithConfig(context.Background(), config())
	if err != nil {
		slog.Error("Err create new config", err)
		return nil, err
	}
	connectiom, err := connPool.Acquire(context.Background())
	if err != nil {
		slog.Error("Err to connect database", err)
		return nil, err
	}
	defer connectiom.Release()

	err = connectiom.Ping(context.Background())
	if err != nil {
		slog.Error("Err ping", err)
		return nil, err
	}
	slog.Info("Connect to db")

	return connPool, nil

}
