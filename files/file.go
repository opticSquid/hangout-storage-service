package files

import (
	"context"
	"errors"
	"regexp"

	"github.com/knadh/koanf/v2"
	"hangout.com/core/storage-service/database"
	"hangout.com/core/storage-service/database/model"
	"hangout.com/core/storage-service/logger"
)

type File struct {
	ContentType string
	Filename    string
	UserId      int32
}

func (f *File) Process(workerId int, ctx context.Context, cfg *koanf.Koanf, dbConnPool *database.DatabaseConnectionPool, log logger.Log) error {
	isVideo, _ := regexp.MatchString(`^video/`, f.ContentType)
	var mediaFile mediaFile
	if isVideo {
		mediaFile = &video{filename: f.Filename}
	} else {
		log.Debug("unsupported content type. can not process file", "contentType", f.ContentType, "file", f.Filename, "worker-id", workerId)
		return errors.New("unsupported file type received, contentType is: " + f.ContentType)
	}
	log.Info("marking file status as PROCESSING in db", "worker-id", workerId, "filename", f.Filename)
	dbConnPool.UpdateProcessingStatus(ctx, f.Filename, model.PROCESSING, log)
	err := mediaFile.processMedia(workerId, cfg, log)
	if err != nil {
		log.Error("marking file status as FAILED in db", "worker-id", workerId, "filename", f.Filename)
		dbConnPool.UpdateProcessingStatus(ctx, f.Filename, model.FAIL, log)
		return err
	}
	log.Info("marking file status as SUCCESS in db", "worker-id", workerId, "filename", f.Filename)
	dbConnPool.UpdateProcessingStatus(ctx, f.Filename, model.SUCCESS, log)
	return nil
}
