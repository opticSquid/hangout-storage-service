package files

import (
	"context"
	"errors"
	"regexp"
	"time"

	"github.com/knadh/koanf/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"hangout.com/core/storage-service/database"
	"hangout.com/core/storage-service/database/model"
	"hangout.com/core/storage-service/logger"
)

type File struct {
	Context     context.Context
	ContentType string
	Filename    string
	UserId      int32
}

func (f *File) Process(workerContext context.Context, cfg *koanf.Koanf, dbConnPool *database.DatabaseConnectionPool, log logger.Log) error {
	initFileProcessMetrics("file_process_duration", "Duration of file processing in seconds")
	start := time.Now()
	tr := otel.Tracer("hangout.storage.file")
	ctx, span := tr.Start(workerContext, "ProcessFile")
	defer func() {
		span.End()
		duration := time.Since(start).Seconds()
		processDuration.Record(ctx, float64(duration))
	}()
	span.SetAttributes(
		attribute.String("file.name", f.Filename),
		attribute.Int("file.userId", int(f.UserId)),
	)
	log = log.With("file", f.Filename, "userId", f.UserId)
	isVideo, _ := regexp.MatchString(`^video/`, f.ContentType)
	var mediaFile mediaFile
	if isVideo {
		mediaFile = &video{filename: f.Filename}
	} else {
		log.Debug(ctx, "unsupported content type. can not process file", "contentType", f.ContentType, "file", f.Filename)
		return errors.New("unsupported file type received, contentType is: " + f.ContentType)
	}
	log.Info(ctx, "marking file status as PROCESSING in db", "filename", f.Filename)
	dbConnPool.UpdateProcessingStatus(ctx, f.Filename, model.PROCESSING, log)
	err := mediaFile.processMedia(ctx, cfg, log)
	if err != nil {
		log.Error(ctx, "marking file status as FAILED in db", "filename", f.Filename)
		dbConnPool.UpdateProcessingStatus(ctx, f.Filename, model.FAIL, log)
		return err
	}
	log.Info(ctx, "marking file status as SUCCESS in db", "filename", f.Filename)
	dbConnPool.UpdateProcessingStatus(ctx, f.Filename, model.SUCCESS, log)
	return nil
}
