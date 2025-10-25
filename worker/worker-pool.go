package worker

import (
	"context"
	"sync"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/knadh/koanf/v2"

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
	s3Client, err := cloudstorage.Connect(workerId, worker.ctx, worker.cfg, workerLogger)
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
			worker.do(workerId, file, workerLogger, s3Client)
		case <-worker.ctx.Done():
			workerLogger.Info(worker.ctx, "Context cancelled, stopping worker")
			return
		}
	}
}

func (worker *WorkerPool) do(workerId int, file *files.File, workerLogger logger.Log, s3Client *s3.Client) {
	tr := otel.Tracer("hangout.storage.worker")
	ctx, span := tr.Start(file.Context, "ProcessFile")
	span.SetAttributes(
		attribute.Int("worker.id", workerId),
		attribute.String("file.name", file.Filename),
		attribute.Int("file.userId", int(file.UserId)),
	)
	defer span.End()

	workerLogger.Info(ctx, "starting file processing", "file-name", file.Filename, "user-id", file.UserId)

	// Check if already processed
	isProcessed, err := worker.dbConnPool.IsAlreadyProcessed(ctx, file.Filename)
	if err != nil {
		workerLogger.Error(ctx, "error checking process status", "error", err.Error())
		// Optionally: do not acknowledge, so it can be retried
		return
	}
	if isProcessed {
		workerLogger.Info(ctx, "file already processed, acknowledging and skipping", "file-name", file.Filename)
		if file.KafkaSession != nil && file.KafkaMessage != nil {
			file.KafkaSession.MarkMessage(file.KafkaMessage, "")
		}
		return
	}

	// Not processed: download, process, upload, then acknowledge
	cloudstorage.Download(ctx, s3Client, file, worker.cfg, workerLogger)
	err = file.Process(ctx, worker.cfg, worker.dbConnPool, workerLogger)
	if err != nil {
		workerLogger.Error(ctx, "could not process file", "error", err.Error())
		// Optionally: do not acknowledge, so it can be retried
		return
	}

	cloudstorage.UploadDir(ctx, s3Client, file, worker.cfg, workerLogger)
	if file.KafkaSession != nil && file.KafkaMessage != nil {
		file.KafkaSession.MarkMessage(file.KafkaMessage, "")
	}
	workerLogger.Info(ctx, "finished file processing", "file-name", file.Filename)
}

// Wait ensures all workers complete processing before the program exits
func (worker *WorkerPool) Wait() {
	worker.wg.Wait()
}
