package h264

import (
	"context"
	"os/exec"

	"hangout.com/core/storage-service/logger"
)

func ProcessSDRResolutions(workerId int, ctx context.Context, inputFilePath string, outputFilePath string, log logger.Log) error {
	log.Info(ctx, "pipeline checkpoint", "file", inputFilePath, "enocder", "h264", "media-type", "video-sdr", "status", "starting processing", "worker-id", workerId)
	var err error
	err = process640p(workerId, ctx, inputFilePath, outputFilePath, log)
	if err != nil {
		return err
	}
	err = process1280p(workerId, ctx, inputFilePath, outputFilePath, log)
	if err != nil {
		return err
	}
	err = process1920p(workerId, ctx, inputFilePath, outputFilePath, log)
	if err != nil {
		return err
	}
	log.Info(ctx, "pipeline checkpoint", "file", inputFilePath, "enocder", "h264", "media-type", "video-sdr", "status", "finished", "worker-id", workerId)
	return nil
}
func process640p(workerId int, ctx context.Context, inputFilePath string, outputFilePath string, log logger.Log) error {
	log.Debug(ctx, "pipeline checkpoint", "file", inputFilePath, "encoder", "h264", "media-type", "video-sdr", "resolution", "360x640", "status", "starting processing", "worker-id", workerId)
	var cmd *exec.Cmd
	var err error
	outputFilePath = outputFilePath + "_h264_640p.mp4"
	cmd = exec.Command("ffmpeg", "-i", inputFilePath, "-c:v", "libx264", "-crf", "25", "-g", "30", "-vf", "scale=320x640", "-preset", "slow", "-an", outputFilePath)
	_, err = cmd.Output()
	if err != nil {
		log.Error(ctx, "error in processing video", "file", inputFilePath, "encoder", "h264", "resolution", "360x640", "error", err.Error(), "worker-id", workerId)
		return err
	} else {
		log.Debug(ctx, "pipeline checkpoint", "file", inputFilePath, "encoder", "h264", "media-type", "video-sdr", "resolution", "360x640", "status", "finished", "worker-id", workerId)
	}
	return nil
}

func process1280p(workerId int, ctx context.Context, inputFilePath string, outputFilePath string, log logger.Log) error {
	log.Debug(ctx, "pipeline checkpoint", "file", inputFilePath, "encoder", "h264", "media-type", "video-sdr", "resolution", "720x1280", "status", "starting processing", "worker-id", workerId)
	var cmd *exec.Cmd
	var err error
	outputFilePath = outputFilePath + "_h264_1280p.mp4"
	cmd = exec.Command("ffmpeg", "-i", inputFilePath, "-c:v", "libx264", "-crf", "25", "-g", "30", "-vf", "scale=720x1280", "-preset", "slow", "-an", outputFilePath)
	_, err = cmd.Output()
	if err != nil {
		log.Error(ctx, "error in processing video", "file", inputFilePath, "encoder", "h264", "resolution", "720x1280", "error", err.Error(), "worker-id", workerId)
		return err
	} else {
		log.Debug(ctx, "pipeline checkpoint", "file", inputFilePath, "encoder", "h264", "media-type", "video-sdr", "resolution", "720x1280", "status", "finished", "worker-id", workerId)
	}
	return nil
}

func process1920p(workerId int, ctx context.Context, inputFilePath string, outputFilePath string, log logger.Log) error {
	log.Debug(ctx, "pipeline checkpoint", "file", inputFilePath, "encoder", "h264", "media-type", "video-sdr", "resolution", "1080x1920", "status", "starting processing", "worker-id", workerId)
	var cmd *exec.Cmd
	var err error
	outputFilePath = outputFilePath + "_h264_1920p.mp4"
	cmd = exec.Command("ffmpeg", "-i", inputFilePath, "-c:v", "libx264", "-crf", "25", "-g", "30", "-vf", "scale=1080x1920", "-preset", "slow", "-an", outputFilePath)
	_, err = cmd.Output()
	if err != nil {
		log.Error(ctx, "error in processing video", "file", inputFilePath, "encoder", "h264", "resolution", "1080x1920", "error", err.Error(), "worker-id", workerId)
		return err
	} else {
		log.Debug(ctx, "pipeline checkpoint", "file", inputFilePath, "encoder", "h264", "media-type", "video-sdr", "resolution", "1080x1920", "status", "finished", "worker-id", workerId)
	}
	return nil
}
