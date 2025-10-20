package worker

import (
	"context"
	"sync"

	"github.com/knadh/koanf/v2"

	"github.com/minio/minio-go/v7"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"hangout.com/core/storage-service/cloudstorage"
	"hangout.com/core/storage-service/database"
	"hangout.com/core/storage-service/files"
	"hangout.com/core/storage-service/logger"
)

type WorkerPool struct {
	eventChan  <-chan *files.File
	wg         *sync.WaitGroup
	ctx        context.Context
	cfg        *koanf.Koanf
	dbConnPool *database.DatabaseConnectionPool
	log        logger.Log
}

func CreateWorkerPool(eventChan <-chan *files.File, ctx context.Context, cfg *koanf.Koanf, dbConnPool *database.DatabaseConnectionPool, log logger.Log) *WorkerPool {
	wp := &WorkerPool{eventChan: eventChan, wg: &sync.WaitGroup{}, ctx: ctx, cfg: cfg, dbConnPool: dbConnPool, log: log}
	for i := 0; i < cfg.Int("process.pool-strength"); i++ {
		log.Debug(ctx, "spawning worker", "worker-id", i)
		wp.wg.Add(1)
		go wp.worker(i)
	}
	return wp
}

func (worker *WorkerPool) worker(workerId int) {
	defer worker.wg.Done()
	workerLogger := worker.log.With("worker-id", workerId)
	minioClient, err := cloudstorage.Connect(workerId, worker.ctx, worker.cfg, workerLogger)
	if err != nil {
		return
	}
	for {
		select {
		case file, ok := <-worker.eventChan:
			if !ok {
				workerLogger.Info(worker.ctx, "Event channel closed, stopping worker")
				return
			}
			worker.do(workerId, file, workerLogger, minioClient)
		case <-worker.ctx.Done():
			workerLogger.Info(worker.ctx, "Context cancelled, stopping worker")
			return
		}
	}
}

func (worker *WorkerPool) do(workerId int, file *files.File, workerLogger logger.Log, minioClient *minio.Client) {
	tr := otel.Tracer("hangout.storage.worker")
	ctx, span := tr.Start(file.Context, "ProcessFile")
	span.SetAttributes(
		attribute.Int("worker.id", workerId),
		attribute.String("file.name", file.Filename),
		attribute.Int("file.userId", int(file.UserId)),
	)
	defer span.End()

	workerLogger.Info(ctx, "starting file processing", "file-name", file.Filename, "user-id", file.UserId)

	// download the given file from cloud storage
	cloudstorage.Download(ctx, minioClient, file, worker.cfg, workerLogger)
	// process the file
	err := file.Process(ctx, worker.cfg, worker.dbConnPool, workerLogger)
	if err != nil {
		workerLogger.Error(ctx, "could not process file", "error", err.Error())
	}
	// upload the given file to cloud storage
	cloudstorage.UploadDir(ctx, minioClient, file, worker.cfg, workerLogger)
	if file.KafkaSession != nil && file.KafkaMessage != nil {
		file.KafkaSession.MarkMessage(file.KafkaMessage, "") // Acknowledge the Kafka event and mark as completed
	}
	workerLogger.Info(ctx, "finished file processing", "file-name", file.Filename)
}

// Wait ensures all workers complete processing before the program exits
func (worker *WorkerPool) Wait() {
	worker.wg.Wait()
}
