package abr

import (
	"context"
	"os/exec"

	"hangout.com/core/storage-service/logger"
)

func CreatePlaylist(workerId int, ctx context.Context, outputFilePath string, encoding string, log logger.Log) error {
	log.Info(ctx, "pipeline status status", "segementation and playlist creation", "starting", "worker-id", workerId)
	videoFile := outputFilePath + "_" + encoding + "_"
	audioFile := outputFilePath + "_" + encoding + "_audio.mp4"
	var cmd *exec.Cmd
	var err error
	log.Debug(ctx, "Input", "video file path", videoFile, "worker-id", workerId)
	log.Debug(ctx, "Input", "audio file path", audioFile, "worker-id", workerId)
	log.Debug(ctx, "Input", "output file path", outputFilePath, "worker-id", workerId)
	cmd = exec.Command("MP4Box", "-dash", "2000", "-frag", "2000", "-segment-name", "segment_$RepresentationID$_", "-fps", "30", videoFile+"640p.mp4#video:id=640p", videoFile+"1280p.mp4#video:id=1280p", videoFile+"1920p.mp4#video:id=1920p", audioFile+"#audio:id=English:role=main", "-out", outputFilePath+".mpd")
	_, err = cmd.Output()
	if err != nil {
		log.Error(ctx, "error in processing segmentation and playlist creation", "error", err.Error(), "worker-id", workerId)
		return err
	} else {
		log.Info(ctx, "pipeline status status", "segementation and playlist creation", "finished", "worker-id", workerId)
	}
	return nil
}
