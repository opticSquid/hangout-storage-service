package vp9

import (
	"context"
	"os/exec"

	"hangout.com/core/storage-service/logger"
)

func ProcessAudio(ctx context.Context, inputFilePath string, outputFilePath string, log logger.Log) error {
	log.Info(ctx, "pipeline checkpoint", "file", inputFilePath, "enocder", "libopus", "media-type", "audio", "status", "starting processing")
	var cmd *exec.Cmd
	var err error
	outputFilePath = outputFilePath + "_vp9_audio.opus"
	log.Debug(ctx, "Input", "Input file path", inputFilePath)
	log.Debug(ctx, "Input", "output file path", outputFilePath)
	cmd = exec.Command("ffmpeg", "-i", inputFilePath, "-vn", "-c:a", "libopus", outputFilePath)
	_, err = cmd.Output()
	if err != nil {
		log.Error(ctx, "error in processing audio", "file", inputFilePath, "encoder", "libopus", "error", err.Error())
		return err
	} else {
		log.Debug(ctx, "pipeline checkpoint", "file", inputFilePath, "encoder", "libopus", "media-type", "audio", "status", "finished")
	}
	return nil
}
