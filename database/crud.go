package database

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"hangout.com/core/storage-service/database/model"
	"hangout.com/core/storage-service/logger"
)

// IsAlreadyProcessed checks if the file with the given filename has process_status SUCCESS.
func (dbConn *DatabaseConnectionPool) IsAlreadyProcessed(ctx context.Context, filename string) (bool, error) {
	var currentStatus model.ProcessStatus
	row := dbConn.pool.QueryRow(ctx, "SELECT process_status FROM media WHERE filename = $1", filename)
	err := row.Scan(&currentStatus)
	if err != nil {
		return false, err
	}
	return currentStatus == model.SUCCESS, nil
}

func (dbConn *DatabaseConnectionPool) UpdateProcessingStatus(ctx context.Context, filename string, processStatus model.ProcessStatus, log logger.Log) error {
	tr := otel.Tracer("hangout.storage.database")
	ctx, span := tr.Start(ctx, "UpdateProcessingStatus")
	span.SetAttributes(
		attribute.String("db.operation", "UPDATE"),
		attribute.String("db.filename", filename),
		attribute.String("db.process_status", string(processStatus)),
	)
	defer span.End()

	query := `UPDATE media SET process_status = $1 where filename = $2`
	cmdTag, err := dbConn.pool.Exec(ctx, query, processStatus, filename)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		log.Error(ctx, "could not update file processing status in database", "error", err)
	}
	if cmdTag.RowsAffected() == 0 {
		log.Error(ctx, "file not found for update", "filename", filename)
	}
	if err != nil {
		return err
	} else {
		return nil
	}
}
