package database

import (
	"context"

	"hangout.com/core/storage-service/database/model"
	"hangout.com/core/storage-service/logger"
)

func (dbConn *DatabaseConnectionPool) UpdateProcessingStatus(ctx context.Context, filename string, processStatus model.ProcessStatus, log logger.Log) error {
	query := `UPDATE media SET process_status = $1 where filename = $2`
	cmdTag, err := dbConn.pool.Exec(ctx, query, processStatus, filename)
	if err != nil {
		log.Error("could not update file processing status in database", "error", err)
	}
	if cmdTag.RowsAffected() == 0 {
		log.Error("file not found for update", "filename", filename)
	}
	if err != nil {
		return err
	} else {
		return nil
	}
}
