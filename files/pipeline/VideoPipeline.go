package pipeline

import (
	"context"
	"os"
	"strings"

	"github.com/knadh/koanf/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"hangout.com/core/storage-service/files/abr"
	"hangout.com/core/storage-service/files/h264"
	"hangout.com/core/storage-service/files/postprocess"
	"hangout.com/core/storage-service/files/vp9"
	"hangout.com/core/storage-service/logger"
)

type Video struct {
	Filename string
}

func (v *Video) ProcessMedia(ctx context.Context, cfg *koanf.Koanf, log logger.Log) error {
	tr := otel.Tracer("hangout.storage.files")
	ctx, span := tr.Start(ctx, "ProcessVideo")
	defer span.End()
	span.SetAttributes(
		attribute.String("video.filename", v.Filename),
	)
	splittedFilename := strings.Split(v.Filename, ".")
	inputFile := "/tmp/" + v.Filename
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
	postprocess.CleanUp(ctx, "h264", v.Filename, log)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func processH264(ctx context.Context, inputFilePath string, outputFolder string, filename string, log logger.Log) error {
	tr := otel.Tracer("hangout.storage.files")
	ctx, span := tr.Start(ctx, "ProcessH264")
	defer span.End()
	span.SetAttributes(
		attribute.String("video.filename", filename),
		attribute.String("encoder", "h264"),
	)
	log = log.With("encoder", "h264")
	log.Info(ctx, "pipeline checkpoint", "status", "starting processing")
	outputFilePath := outputFolder + "/" + filename
	log.Debug(ctx, "Check", "Input file path", inputFilePath)
	log.Debug(ctx, "Check", "Output file path", outputFilePath)
	h264.ProcessSDRResolutions(ctx, inputFilePath, outputFilePath, log)
	h264.ProcessAudio(ctx, inputFilePath, outputFilePath, log)
	abr.CreatePlaylist(ctx, outputFilePath, "h264", log)
	log.Info(ctx, "pipeline checkpoint", "file", inputFilePath, "status", "finished processing")
	return nil
}

func processVp9(ctx context.Context, inputFilePath string, outputFolder string, filename string, log logger.Log) error {
	tr := otel.Tracer("hangout.storage.files")
	ctx, span := tr.Start(ctx, "ProcessVp9")
	defer span.End()
	span.SetAttributes(
		attribute.String("video.filename", filename),
		attribute.String("encoder", "vp9"),
	)
	log = log.With("encoder", "vp9")
	log.Info(ctx, "pipeline checkpoint", "status", "starting processing")
	outputFilePath := outputFolder + "/" + filename
	log.Debug(ctx, "Input")
	log.Debug(ctx, "Output", "output file path", outputFilePath)
	vp9.ProcessSDRResolutions(ctx, inputFilePath, outputFilePath, log)
	vp9.ProcessAudio(ctx, inputFilePath, outputFilePath, log)
	abr.CreatePlaylist(ctx, outputFilePath, "vp9", log)
	log.Info(ctx, "pipeline checkpoint", "status", "finished processing")
	return nil
}
