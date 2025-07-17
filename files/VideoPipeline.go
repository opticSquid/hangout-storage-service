package files

import (
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

func (v *video) processMedia(workerId int, cfg *koanf.Koanf, log logger.Log) error {
	splittedFilename := strings.Split(v.filename, ".")
	inputFile := "/tmp/" + v.filename
	outputFolder := "/tmp/" + splittedFilename[0]
	filename := splittedFilename[0]
	var err error
	err = os.Mkdir(outputFolder, 0755)
	if err != nil {
		log.Error("could not create base output folder", "err", err.Error(), "worker-id", workerId)
	}
	err = processH264(workerId, inputFile, outputFolder, filename, log)
	if err != nil {
		log.Error("error in video processing pipeline", "error", err.Error(), "worker-id", workerId)
	}
	postprocess.CleanUp(workerId, "h264", v.filename, log)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func processH264(workerId int, inputFilePath string, outputFolder string, filename string, log logger.Log) error {
	log.Info("pipeline checkpoint", "file", inputFilePath, "enocder", "h264", "status", "starting processing", "worker-id", workerId)
	outputFilePath := outputFolder + "/" + filename
	log.Debug("Check", "Input file path", inputFilePath, "worker-id", workerId)
	log.Debug("Check", "Output file path", outputFilePath, "worker-id", workerId)
	h264.ProcessSDRResolutions(workerId, inputFilePath, outputFilePath, log)
	h264.ProcessAudio(workerId, inputFilePath, outputFilePath, log)
	abr.CreatePlaylist(workerId, outputFilePath, "h264", log)
	log.Info("pipeline checkpoint", "file", inputFilePath, "enocder", "h264", "status", "finished processing", "worker-id", workerId)
	return nil
}

func processVp9(workerId int, inputFilePath string, outputFolder string, filename string, log logger.Log) error {
	log.Info("pipeline checkpoint", "file", inputFilePath, "enocder", "vp9", "status", "starting processing")
	outputFilePath := outputFolder + "/" + filename
	log.Debug("Input", "Input file path", inputFilePath)
	log.Debug("Input", "output file path", outputFilePath)
	vp9.ProcessSDRResolutions(inputFilePath, outputFilePath, log)
	vp9.ProcessAudio(inputFilePath, outputFilePath, log)
	abr.CreatePlaylist(workerId, outputFilePath, "vp9", log)
	log.Info("pipeline checkpoint", "file", inputFilePath, "enocder", "vp9", "status", "finished processing")
	return nil
}
