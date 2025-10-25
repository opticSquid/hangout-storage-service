package files

import (
	"context"
	"errors"
	"regexp"
	"time"

	"github.com/IBM/sarama"
	"github.com/knadh/koanf/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"hangout.com/core/storage-service/database"
	"hangout.com/core/storage-service/database/model"
	"hangout.com/core/storage-service/files/pipeline"
	"hangout.com/core/storage-service/logger"
)

type File struct {
	Context      context.Context
	ContentType  string
	Filename     string
	UserId       int32
	KafkaMessage *sarama.ConsumerMessage
	KafkaSession sarama.ConsumerGroupSession
}

func (f *File) Process(workerContext context.Context, cfg *koanf.Koanf, dbConnPool *database.DatabaseConnectionPool, log logger.Log) error {
	start := time.Now()
	initFileProcessMetrics("file_process_duration", "Duration of file processing in seconds")
	tr := otel.Tracer("hangout.storage.file")
	ctx, span := tr.Start(workerContext, "ProcessFile")
	defer func() {
		span.End()
		duration := time.Since(start).Minutes()
		processDuration.Record(ctx, float64(duration))
	}()
	span.SetAttributes(
		attribute.String("file.name", f.Filename),
		attribute.Int("file.userId", int(f.UserId)),
		attribute.String("file.contentType", f.ContentType),
	)
	log = log.With("file", f.Filename, "userId", f.UserId)
	isVideo, _ := regexp.MatchString(`^video/`, f.ContentType)
	if !isVideo {
		log.Debug(ctx, "unsupported content type. can not process file", "contentType", f.ContentType, "file", f.Filename)
		span.SetStatus(codes.Error, "Unsupported content type")
		span.RecordError(errors.New("unsupported content type"))
		return errors.New("unsupported file type received, contentType is: " + f.ContentType)
	} else {
		mediaFile := &pipeline.Video{Filename: f.Filename}
		log.Info(ctx, "marking file status as PROCESSING in db", "filename", f.Filename)
		err := dbConnPool.UpdateProcessingStatus(ctx, f.Filename, model.PROCESSING, log)
		if err != nil {
			log.Error(ctx, "could not mark file as PROCESSING in db", "filename", f.Filename)
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return err
		}
		err = mediaFile.ProcessMedia(ctx, cfg, log)
		if err != nil {
			log.Error(ctx, "marking file status as FAILED in db", "filename", f.Filename)
			dbConnPool.UpdateProcessingStatus(ctx, f.Filename, model.FAIL, log)
			return err
		}
		log.Info(ctx, "marking file status as SUCCESS in db", "filename", f.Filename)
		dbConnPool.UpdateProcessingStatus(ctx, f.Filename, model.SUCCESS, log)
		return nil
	}
}
