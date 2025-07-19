package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/knadh/koanf/v2"
	"hangout.com/core/storage-service/config"
	"hangout.com/core/storage-service/database"
	"hangout.com/core/storage-service/files"
	"hangout.com/core/storage-service/kafka"
	"hangout.com/core/storage-service/logger"
	"hangout.com/core/storage-service/telemetry"
	"hangout.com/core/storage-service/worker"
)

var CONFIG = koanf.New(".")

func main() {
	config.InitAppConfig(CONFIG)
	log := logger.NewLogger(CONFIG)
	// Create a base context with a cancel function for the application lifecycle
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Ensure cancel is called on application exit

	// Handle OS signals for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Info(ctx, "Received shutdown signal, cancelling context")
		cancel()
	}()

	log.Info(ctx, "starting Hangout Storage Service")

	// Initialize Open Telemetry sdk
	log.Info(ctx, "setting up telemetry")
	otelShutdown, err := telemetry.SetUpOTelSDK(ctx, CONFIG, log)
	if err != nil {
		log.Error(ctx, "could not set up telemetry", "error", err)
		return
	}
	// Handle shutdown properly so nothing leaks.
	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()
	// Start the database connection
	dbConnpool := database.ConnectToDB(ctx, CONFIG, log)
	defer dbConnpool.Close(ctx, log)

	// Channel to handle incoming Kafka events
	eventChan := make(chan *files.File, CONFIG.Int("process.queue-length"))

	// Start the worker pool with the base context
	log.Info(ctx, "Creating worker pool", "pool-strength", CONFIG.Int("process.queue-length"))
	wp := worker.CreateWorkerPool(eventChan, ctx, CONFIG, dbConnpool, log)

	// Start the Kafka consumer
	err = kafka.StartConsumer(eventChan, ctx, CONFIG, log)
	if err != nil {
		log.Error(ctx, "Error starting Consumer Group")
	}

	// Wait for all workers to finish on shutdown
	wp.Wait()
	log.Info(ctx, "Hangout Storage Service shut down gracefully")
}
