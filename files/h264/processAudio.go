package h264

import (
	"context"
	"os/exec"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"hangout.com/core/storage-service/logger"
)

func ProcessAudio(ctx context.Context, inputFilePath string, outputFilePath string, log logger.Log) error {
	tr := otel.Tracer("hangout.storage.video")
	ctx, span := tr.Start(ctx, "ProcessH264Audio")
	defer span.End()
	span.SetAttributes(
		attribute.String("video.filename", outputFilePath),
		attribute.String("media-type", "audio"),
		attribute.String("encoder", "aac"),
	)
	log = log.With("encoder", "aac")
	log.Info(ctx, "pipeline checkpoint", "status", "starting processing")
	var cmd *exec.Cmd
	var err error
	outputFilePath = outputFilePath + "_h264_audio.mp4"
	log.Debug(ctx, "Input", "Input file path", inputFilePath)
	log.Debug(ctx, "Input", "output file path", outputFilePath)
	cmd = exec.Command("ffmpeg", "-i", inputFilePath, "-vn", "-c:a", "aac", outputFilePath)
	_, err = cmd.Output()
	if err != nil {
		log.Error(ctx, "error in processing audio", "error", err.Error())
		return err
	} else {
		log.Debug(ctx, "pipeline checkpoint", "status", "finished")
	}
	return nil
}
