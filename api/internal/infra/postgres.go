package infra

import (
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"log/slog"
)

func InitPostgres(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Open("pgx", dsn)

	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	slog.Info("connected to postgres")

	return db, err
}
