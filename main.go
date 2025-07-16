package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/knadh/koanf/v2"
	"hangout.com/core/storage-service/config"
	"hangout.com/core/storage-service/files"
	"hangout.com/core/storage-service/kafka"
	"hangout.com/core/storage-service/logger"
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
		log.Info("Received shutdown signal, cancelling context")
		cancel()
	}()

	log.Info("starting Hangout Storage Service", "logging-backend", CONFIG.String("log.backend"))

	// Channel to handle incoming Kafka events
	eventChan := make(chan *files.File, CONFIG.Int("process.queue-length"))

	// Start the worker pool with the base context
	log.Info("Creating worker pool", "pool-strength", CONFIG.Int("process.queue-length"))
	wp := worker.CreateWorkerPool(eventChan, ctx, CONFIG, log)

	// Start the Kafka consumer
	log.Info("starting kafka consumer using ConsumerGroup API")
	err := kafka.StartConsumer(eventChan, ctx, CONFIG, log)
	if err != nil {
		log.Error("Error starting Consumer Group")
	}

	// Wait for all workers to finish on shutdown
	wp.Wait()
	log.Info("Hangout Storage Service shut down gracefully")
}
