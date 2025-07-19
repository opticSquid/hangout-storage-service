package files

import (
	"context"
	"os"
	"strings"

	"github.com/knadh/koanf/v2"
	"hangout.com/core/storage-service/files/abr"
	"hangout.com/core/storage-service/files/h264"
	"hangout.com/core/storage-service/files/postprocess"
	"hangout.com/core/storage-service/files/vp9"
	"hangout.com/core/storage-service/logger"
)

type video struct {
	filename string
}

func (v *video) processMedia(ctx context.Context, cfg *koanf.Koanf, log logger.Log) error {
	splittedFilename := strings.Split(v.filename, ".")
	inputFile := "/tmp/" + v.filename
	outputFolder := "/tmp/" + splittedFilename[0]
	filename := splittedFilename[0]
	var err error
	err = os.Mkdir(outputFolder, 0755)
	if err != nil {
		log.Error(ctx, "could not create base output folder", "err", err.Error())
	}
	err = processH264(ctx, inputFile, outputFolder, filename, log)
	if err != nil {
		log.Error(ctx, "error in video processing pipeline", "error", err.Error())
	}
	postprocess.CleanUp(ctx, "h264", v.filename, log)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func processH264(ctx context.Context, inputFilePath string, outputFolder string, filename string, log logger.Log) error {
	log.Info(ctx, "pipeline checkpoint", "file", inputFilePath, "enocder", "h264", "status", "starting processing")
	outputFilePath := outputFolder + "/" + filename
	log.Debug(ctx, "Check", "Input file path", inputFilePath)
	log.Debug(ctx, "Check", "Output file path", outputFilePath)
	h264.ProcessSDRResolutions(ctx, inputFilePath, outputFilePath, log)
	h264.ProcessAudio(ctx, inputFilePath, outputFilePath, log)
	abr.CreatePlaylist(ctx, outputFilePath, "h264", log)
	log.Info(ctx, "pipeline checkpoint", "file", inputFilePath, "enocder", "h264", "status", "finished processing")
	return nil
}

func processVp9(ctx context.Context, inputFilePath string, outputFolder string, filename string, log logger.Log) error {
	log.Info(ctx, "pipeline checkpoint", "file", inputFilePath, "enocder", "vp9", "status", "starting processing")
	outputFilePath := outputFolder + "/" + filename
	log.Debug(ctx, "Input", "Input file path", inputFilePath)
	log.Debug(ctx, "Input", "output file path", outputFilePath)
	vp9.ProcessSDRResolutions(ctx, inputFilePath, outputFilePath, log)
	vp9.ProcessAudio(ctx, inputFilePath, outputFilePath, log)
	abr.CreatePlaylist(ctx, outputFilePath, "vp9", log)
	log.Info(ctx, "pipeline checkpoint", "file", inputFilePath, "enocder", "vp9", "status", "finished processing")
	return nil
}
