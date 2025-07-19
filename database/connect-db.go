package database

import (
	"context"
	"fmt"
	"net/url"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/knadh/koanf/v2"
	"hangout.com/core/storage-service/exceptions"
	"hangout.com/core/storage-service/logger"
)

type DatabaseConnectionPool struct {
	pool *pgxpool.Pool
}

func ConnectToDB(ctx context.Context, cfg *koanf.Koanf, log logger.Log) *DatabaseConnectionPool {
	log.Info(ctx, "connecting to database")
	password := cfg.String("datasource.password")
	password = url.QueryEscape(password)
	dbConnectionString := fmt.Sprintf("postgres://%s:%s@%s/%s", cfg.String("datasource.username"), password, cfg.String("datasource.url"), cfg.String("datasource.dbname"))
	dbConnPool, err := pgxpool.New(ctx, dbConnectionString)
	if err != nil {
		exceptions.DbConnectionError(ctx, "could not connect to database", &err, log)
	}
	log.Info(ctx, "successfully connected to database")
	return &DatabaseConnectionPool{pool: dbConnPool}
}

func (dbConn *DatabaseConnectionPool) Close(ctx context.Context, log logger.Log) {
	if dbConn.pool != nil {
		dbConn.pool.Close()
		log.Info(ctx, "closed database connection")
	}
}
