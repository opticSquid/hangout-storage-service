package files

import (
	"os"
	"os/exec"
	"strings"

	"github.com/knadh/koanf/v2"
	"hangout.com/core/storage-service/logger"
)

type image struct {
	filename string
}

func (i *image) processMedia(workerId int, cfg *koanf.Koanf, log logger.Log) error {
	splittedFilename := strings.Split(i.filename, ".")
	inputFile := "/tmp/" + i.filename
	outputFolder := "/tmp/" + splittedFilename[0]
	var err error
	err = os.Mkdir(outputFolder, 0755)
	if err != nil {
		log.Error("could not create base output folder", "err", err.Error())
		panic(err)
	}
	outputFile := outputFolder + "/" + splittedFilename[0]
	err = convertToAvif(inputFile, outputFile, log)
	if err != nil {
		log.Error("error in avif pipeline", "error", err.Error())
	}
	return nil
}

func convertToAvif(inputFile string, outputFile string, log logger.Log) error {
	log.Info("starting to convert Image to AVIF")
	log.Debug("resizing image to 144px width")
	cmd := exec.Command("convert", inputFile, "-resize", "144", outputFile+"_144.avif")
	_, err := cmd.Output()
	if err != nil {
		log.Error("error in converting image to avif")
	}
	log.Debug("resizing image to 360px width")
	cmd = exec.Command("convert", inputFile, "-resize", "360", outputFile+"_360.avif")
	_, err = cmd.Output()
	if err != nil {
		log.Error("error in converting image to avif")
	}
	log.Debug("resizing image to 720px width")
	cmd = exec.Command("convert", inputFile, "-resize", "720", outputFile+"_720.avif")
	_, err = cmd.Output()
	if err != nil {
		log.Error("error in converting image to avif")
	}
	log.Debug("resizing image to 1080px width")
	cmd = exec.Command("convert", inputFile, "-resize", "1080", outputFile+"_1080.avif")
	_, err = cmd.Output()
	if err != nil {
		log.Error("error in converting image to avif")
	}
	return err
}
