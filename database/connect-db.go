package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/knadh/koanf/v2"
	"hangout.com/core/storage-service/exceptions"
	"hangout.com/core/storage-service/logger"
)

func ConnectToDB(ctx context.Context, cfg *koanf.Koanf, log logger.Log) *pgxpool.Pool {
	log.Info("connecting to database")
	dbpool, err := pgxpool.New(ctx, fmt.Sprintf("postgres://%s:%s@%s/%s", cfg.String("datasource.username"), cfg.String("datasource.password"), cfg.String("datasource.url"), cfg.String("datasource.dbname")))
	if err != nil {
		exceptions.DbConnectionError("could not connect to database", &err, log)
	}
	log.Info("successfully connected to database")
	return dbpool
}
