package postprocess

import (
	"context"
	"os"
	"strings"

	"hangout.com/core/storage-service/logger"
)

func CleanUp(ctx context.Context, encoding string, filename string, log logger.Log) {

	// delete the original file from temp directory
	storageDir := "/tmp"
	sourceFilepath := storageDir + "/" + filename
	err := os.Remove(sourceFilepath)
	if err != nil {
		log.Debug(ctx, "could not delete the original file", "error", err, "path", sourceFilepath)
	}
	log.Debug(ctx, "removed source file", "path", sourceFilepath)

	// remove transcoded files
	baseFilename := strings.Split(filename, ".")[0]
	transcodedVideoFileBaseName := storageDir + "/" + baseFilename + "/" + baseFilename + "_" + encoding + "_"
	resolutions := []string{"640p", "1280p", "1920p", "audio"}
	for _, res := range resolutions {
		finalFileName := transcodedVideoFileBaseName + res + ".mp4"
		log.Debug(ctx, "removing transcoded video files", "file", finalFileName)
		err = os.Remove(finalFileName)
		if err != nil {
			log.Error(ctx, "could not remove transcoded video file", "error", err, "file", finalFileName)
		}
	}
}
