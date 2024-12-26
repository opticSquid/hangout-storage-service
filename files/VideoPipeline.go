package files

import (
	"os"
	"os/exec"
	"strings"

	"hangout.com/core/storage-service/config"
	"hangout.com/core/storage-service/files/abr"
	"hangout.com/core/storage-service/files/h264"
	"hangout.com/core/storage-service/files/vp9"
	"hangout.com/core/storage-service/logger"
)

type video struct {
	filename string
}

func (v *video) processMedia(cfg *config.Config, log logger.Log) error {
	splittedFilename := strings.Split(v.filename, ".")
	inputFile := cfg.Hangout.Media.UploadPath + "/" + v.filename
	outputFolder := cfg.Hangout.Media.ProcessedPath + "/" + splittedFilename[0]
	filename := splittedFilename[0]
	var err error
	err = os.Mkdir(outputFolder, 0755)
	if err != nil {
		log.Error("could not create base output folder", "err", err.Error())
		panic(err)
	}
	outputFolder = outputFolder + "/" + "vp9"
	err = os.Mkdir(outputFolder, 0755)
	if err != nil {
		log.Error("could not create vp9 ouput folder", "err", err.Error())
		panic(err)
	}
	err = processVp9(inputFile, outputFolder, filename, log)
	if err != nil {
		log.Error("error in vp9 pipeline", "error", err.Error())
		panic(err)
	}
	return nil
}

// generte h264 (crf mode and 2 pass mode) encoded versions of the uploaded file in 360p, 720p, 1080p
func processH264(inputFile string, outputFolder string, outputFilename string, log logger.Log) error {
	log.Info("h264 pipeline", "status", "starting")
	outputFile := outputFolder + "/" + outputFilename
	h264.ProcessCRFmode(inputFile, outputFile, log)
	h264.Process2PassMode(inputFile, outputFile, log)
	h264.ProcessAudio(inputFile, outputFile, log)
	abr.CreatePlaylist(outputFile, "h264", log)
	log.Info("h264 pipeline", "status", "finished")
	return nil
}

func processVp9(inputFilePath string, outputFolder string, filename string, log logger.Log) error {
	log.Info("pipeline checkpoint", "file", inputFilePath, "enocder", "vp9", "status", "starting processing")
	outputFilePath := outputFolder + "/" + filename
	vp9.ProcessSDRResolutions(inputFilePath, outputFilePath, log)
	vp9.ProcessAudio(inputFilePath, outputFilePath, log)
	abr.CreatePlaylist(outputFilePath, "vp9", log)
	log.Info("pipeline checkpoint", "file", inputFilePath, "enocder", "vp9", "status", "finished processing")
	return nil
}

// generte h265 (crf mode and 2 pass mode) encoded versions of the uploaded file in 360p, 720p, 1080p
func processH265(inputFile string, filename string, log logger.Log) error {
	log.Info("h265 pipeline", "status", "starting")
	var cmd *exec.Cmd
	var err error

	log.Debug("pipeline status", "encoder", "h265", "method", "crf", "status", "starting")
	// creating 360p, 720p using h265 crf mode
	cmd = exec.Command("ffmpeg", "-i", inputFile,
		"-c:v", "libx265", "-crf", "28", "-maxrate", "500k", "-bufsize", "1M", "-vf", "scale=-2:360", filename+"_h265_360p.mp4",
		"-c:v", "libx265", "-crf", "28", "-maxrate", "1M", "-bufsize", "3M", "-vf", "scale=-2:720", filename+"_h265_720p.mp4",
	)
	_, err = cmd.Output()
	if err != nil {
		log.Error("error in processing h265 crf workflow", "error", err.Error())
	} else {
		log.Debug("pipeline status", "encoder", "h265", "method", "crf", "status", "finished")
	}

	// creating 1080p and original using h265 2 pass mode
	log.Debug("pipeline status", "encoder", "h265", "method", "2 pass", "status", "starting")

	// doing 1st pass for 1080p
	cmd = exec.Command("ffmpeg", "-i", inputFile,
		"-c:v", "libx265", "-x265-params", "pass=1:log-level=2:stats="+filename, "-fps_mode", "cfr", "-b:v", "2M", "-vf", "scale=-2:1080", "-an", "-f", "null", "/dev/null",
	)
	_, err = cmd.Output()
	if err != nil {
		log.Error("error in processing h265 2 pass workflow, error in 1st pass", "current resolution", "1080p", "error", err.Error())
	} else {
		log.Debug("pipeline status", "encoder", "h265", "method", "2 pass", "current pass", 1, "current resolution", "1080p", "status", "finished")
	}
	// doing 2nd pass
	// creating 1080p video in 2nd pass out of 1st pass log and cutree files
	cmd = exec.Command("ffmpeg", "-i", inputFile,
		"-c:v", "libx265", "-x265-params", "pass=2:log-level=2:stats="+filename, "-fps_mode", "cfr", "-b:v", "2M", "-vf", "scale=-2:1080", filename+"_h265_1080p.mp4",
	)
	_, err = cmd.Output()
	if err != nil {
		log.Error("error in processing h265 2 pass workflow, error in 2nd pass", "current resolution", "1080p", "error", err.Error())
	} else {
		log.Debug("pipeline status", "encoder", "h265", "method", "2 pass", "current pass", 2, "current resolution", "1080p", "status", "finished")
	}
	// deleting ffmpeg generated log file for 1080p
	err = os.Remove(filename)
	if err != nil {
		log.Error("error in deleting ffmpeg  h265 log file", "current resolution", "1080p", "error", err.Error())
		return err
	} else {
		log.Debug("deleted ffmpeg h265 log file", "current resolution", "1080p")
	}
	// deleting ffmpeg generated cutree file for 1080p
	err = os.Remove(filename + ".cutree")
	if err != nil {
		log.Error("error in deleting ffmpeg h265 cutree file", "current resolution", "1080p", "error", err.Error())
		return err
	} else {
		log.Debug("deleted ffmpeg h265 cutree file", "current resolution", "1080p")
	}
	log.Debug("pipeline status", "encoder", "h265", "method", "2 pass", "status", "finished")

	log.Info("h265 pipeline", "status", "finished")
	return nil
}