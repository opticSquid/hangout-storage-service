package worker

import (
	"context"
	"sync"

	"github.com/knadh/koanf/v2"
	cloudstorage "hangout.com/core/storage-service/cloudStorage"
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

func (wp *WorkerPool) worker(workerId int) {
	defer wp.wg.Done()
	minioClient, err := cloudstorage.Connect(workerId, wp.ctx, wp.cfg, wp.log)
	if err != nil {
		return
	}
	for {
		select {
		case file, ok := <-wp.eventChan:
			if !ok {
				wp.log.Info(wp.ctx, "Event channel closed, stopping worker", "worker-id", workerId)
				return
			}
			wp.log.Info(wp.ctx, "starting file processing", "file-name", file.Filename, "user-id", file.UserId, "worker-id", workerId)

			// download the given file from cloud storage
			cloudstorage.Download(workerId, wp.ctx, minioClient, file, wp.cfg, wp.log)
			// process the file
			err := file.Process(workerId, wp.ctx, wp.cfg, wp.dbConnPool, wp.log)
			if err != nil {
				wp.log.Error(wp.ctx, "could not process file", "error", err.Error(), "worker-id", workerId)
			}
			// upload the given file to cloud storage
			cloudstorage.UploadDir(workerId, wp.ctx, minioClient, file, wp.cfg, wp.log)
			wp.log.Info(wp.ctx, "finished file processing", "file-name", file.Filename, "worker-id", workerId)
		case <-wp.ctx.Done():
			wp.log.Info(wp.ctx, "Context cancelled, stopping worker", "worker-id", workerId)
			return
		}
	}
}

// Wait ensures all workers complete processing before the program exits
func (wp *WorkerPool) Wait() {
	wp.wg.Wait()
}
