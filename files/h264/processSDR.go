package h264

import (
	"context"
	"os/exec"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"hangout.com/core/storage-service/logger"
)

func ProcessSDRResolutions(ctx context.Context, inputFilePath string, outputFilePath string, log logger.Log) error {
	tr := otel.Tracer("hangout.storage.files.h264")
	ctx, span := tr.Start(ctx, "ProcessSDRResolutions")
	defer span.End()
	span.SetAttributes(
		attribute.String("video.filename", inputFilePath),
		attribute.String("encoder", "h264"),
		attribute.String("media-type", "video-sdr"),
	)
	log.Info(ctx, "pipeline checkpoint", "status", "starting processing")
	var err error
	err = process640p(ctx, inputFilePath, outputFilePath, log)
	if err != nil {
		return err
	}
	err = process1280p(ctx, inputFilePath, outputFilePath, log)
	if err != nil {
		return err
	}
	err = process1920p(ctx, inputFilePath, outputFilePath, log)
	if err != nil {
		return err
	}
	log.Info(ctx, "pipeline checkpoint", "status", "finished")
	return nil
}
func process640p(ctx context.Context, inputFilePath string, outputFilePath string, log logger.Log) error {
	tr := otel.Tracer("hangout.storage.files.h264")
	ctx, span := tr.Start(ctx, "Process640p")
	defer span.End()
	span.SetAttributes(
		attribute.String("video.filename", inputFilePath),
		attribute.String("encoder", "h264"),
		attribute.String("media-type", "video-sdr"),
		attribute.String("resolution", "320x640"),
	)
	log = log.With("resolution", "320x640")
	log.Debug(ctx, "pipeline checkpoint", "status", "starting processing")
	var cmd *exec.Cmd
	var err error
	outputFilePath = outputFilePath + "_h264_640p.mp4"
	cmd = exec.Command("ffmpeg", "-i", inputFilePath, "-c:v", "libx264", "-crf", "25", "-g", "30", "-vf", "scale=320x640", "-preset", "slow", "-an", outputFilePath)
	_, err = cmd.Output()
	if err != nil {
		log.Error(ctx, "error in processing video", "error", err.Error())
		return err
	} else {
		log.Debug(ctx, "pipeline checkpoint", "status", "finished")
	}
	return nil
}

func process1280p(ctx context.Context, inputFilePath string, outputFilePath string, log logger.Log) error {
	tr := otel.Tracer("hangout.storage.files.h264")
	ctx, span := tr.Start(ctx, "Process1280p")
	defer span.End()
	span.SetAttributes(
		attribute.String("video.filename", inputFilePath),
		attribute.String("encoder", "h264"),
		attribute.String("media-type", "video-sdr"),
		attribute.String("resolution", "720x1280"),
	)
	log = log.With("resolution", "720x1280")
	log.Debug(ctx, "pipeline checkpoint", "status", "starting processing")
	var cmd *exec.Cmd
	var err error
	outputFilePath = outputFilePath + "_h264_1280p.mp4"
	cmd = exec.Command("ffmpeg", "-i", inputFilePath, "-c:v", "libx264", "-crf", "25", "-g", "30", "-vf", "scale=720x1280", "-preset", "slow", "-an", outputFilePath)
	_, err = cmd.Output()
	if err != nil {
		log.Error(ctx, "error in processing video", "error", err.Error())
		return err
	} else {
		log.Debug(ctx, "pipeline checkpoint", "status", "finished")
	}
	return nil
}

func process1920p(ctx context.Context, inputFilePath string, outputFilePath string, log logger.Log) error {
	tr := otel.Tracer("hangout.storage.files.h264")
	ctx, span := tr.Start(ctx, "Process1920p")
	defer span.End()
	span.SetAttributes(
		attribute.String("video.filename", inputFilePath),
		attribute.String("encoder", "h264"),
		attribute.String("media-type", "video-sdr"),
		attribute.String("resolution", "1080x1920"),
	)
	log = log.With("resolution", "1080x1920")
	log.Debug(ctx, "pipeline checkpoint", "status", "starting processing")
	var cmd *exec.Cmd
	var err error
	outputFilePath = outputFilePath + "_h264_1920p.mp4"
	cmd = exec.Command("ffmpeg", "-i", inputFilePath, "-c:v", "libx264", "-crf", "25", "-g", "30", "-vf", "scale=1080x1920", "-preset", "slow", "-an", outputFilePath)
	_, err = cmd.Output()
	if err != nil {
		log.Error(ctx, "error in processing video", "error", err.Error())
		return err
	} else {
		log.Debug(ctx, "pipeline checkpoint", "status", "finished")
	}
	return nil
}
