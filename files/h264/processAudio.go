package h264

import (
	"context"
	"os/exec"

	"hangout.com/core/storage-service/logger"
)

func ProcessAudio(ctx context.Context, inputFilePath string, outputFilePath string, log logger.Log) error {
	log.Info(ctx, "pipeline checkpoint", "file", inputFilePath, "enocder", "aac", "media-type", "audio", "status", "starting processing")
	var cmd *exec.Cmd
	var err error
	outputFilePath = outputFilePath + "_h264_audio.mp4"
	log.Debug(ctx, "Input", "Input file path", inputFilePath)
	log.Debug(ctx, "Input", "output file path", outputFilePath)
	cmd = exec.Command("ffmpeg", "-i", inputFilePath, "-vn", "-c:a", "aac", outputFilePath)
	_, err = cmd.Output()
	if err != nil {
		log.Error(ctx, "error in processing audio", "file", inputFilePath, "encoder", "aac", "error", err.Error())
		return err
	} else {
		log.Debug(ctx, "pipeline checkpoint", "file", inputFilePath, "encoder", "aac", "media-type", "audio", "status", "finished")
	}
	return nil
}
