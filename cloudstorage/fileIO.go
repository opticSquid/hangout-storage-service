package cloudstorage

import (
	"context"
	"mime"
	"os"
	"path/filepath"
	"strings"

	"github.com/knadh/koanf/v2"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"hangout.com/core/storage-service/files"
	"hangout.com/core/storage-service/logger"
)

func Connect(workerId int, ctx context.Context, cfg *koanf.Koanf, log logger.Log) (*minio.Client, error) {
	log.Info(ctx, "connecting to Minio/s3", "worker-id", workerId)
	useSsl := false
	minioClient, err := minio.New(cfg.String("minio.base-url"), &minio.Options{Creds: credentials.NewStaticV4(cfg.String("minio.access-key"), cfg.String("minio.secret-key"), ""), Secure: useSsl})
	if err != nil {
		log.Error(ctx, "could not connect to Minio/s3", "url", cfg.String("minio.base-url"), "error", err, "worker-id", workerId)
	}
	log.Info(ctx, "Checking if buckets exist", "worker-id", workerId)
	_, err = minioClient.BucketExists(ctx, cfg.String("minio.upload-bucket"))
	if err != nil {
		log.Info(ctx, "Upload bucket does not exist, Creating a new one", "worker-id", workerId)
		err = minioClient.MakeBucket(ctx, cfg.String("minio.upload-bucket"), minio.MakeBucketOptions{})
		if err != nil {
			log.Error(ctx, "Error in creating new upload bucket", "worker-id", workerId)
		} else {
			log.Info(ctx, "Upload bucket successfully created", "worker-id", workerId)
		}
	} else {
		log.Debug(ctx, "Upload bucket exists skipping creation", "worker-id", workerId)
	}
	_, err = minioClient.BucketExists(ctx, cfg.String("minio.storage-bucket"))
	if err != nil {
		log.Info(ctx, "Storage bucket does not exist, Creating a new one", "worker-id", workerId)
		err = minioClient.MakeBucket(ctx, cfg.String("minio.storage-bucket"), minio.MakeBucketOptions{})
		if err != nil {
			log.Error(ctx, "Error in creating new storage bucket", "worker-id", workerId)
		} else {
			log.Info(ctx, "Storage bucket successfully created", "worker-id", workerId)
		}
	} else {
		log.Debug(ctx, "Storage bucket exists skipping creation", "worker-id", workerId)
	}
	return minioClient, err
}
func Download(ctx context.Context, minioClient *minio.Client, file *files.File, cfg *koanf.Koanf, log logger.Log) {
	tr := otel.Tracer("hangout.storage.cloudstorage")
	ctx, span := tr.Start(ctx, "DownloadFile")
	defer span.End()
	span.SetAttributes(
		attribute.String("file.name", file.Filename),
		attribute.Int("file.userId", int(file.UserId)),
	)
	log.Info(ctx, "Downloading file", "file", file.Filename)
	err := minioClient.FGetObject(ctx, cfg.String("minio.upload-bucket"), file.Filename, "/tmp/"+file.Filename, minio.GetObjectOptions{})
	if err != nil {
		log.Error(ctx, "Error occured while downloading file", "file", file.Filename)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return
	}
}

func UploadDir(ctx context.Context, minioClient *minio.Client, event *files.File, cfg *koanf.Koanf, log logger.Log) {
	tr := otel.Tracer("hangout.storage.cloudstorage")
	ctx, span := tr.Start(ctx, "UploadDir")
	defer span.End()
	span.SetAttributes(
		attribute.String("file.name", event.Filename),
		attribute.Int("file.userId", int(event.UserId)),
	)
	baseFilename := strings.Split(event.Filename, ".")[0]
	currentDir := "/tmp/" + baseFilename
	log.Info(ctx, "Starting to upload directory to Minio/s3", "directory", currentDir)
	// Walk through the folder and upload files
	err := filepath.Walk(currentDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Error(ctx, "Error occured while traversing through current file in the directory", "file", event.Filename, "error", err)
			return err
		}
		log.Debug(ctx, "trying to upload file", "file name", info.Name(), "path", path)
		// Skip directories
		if info.IsDir() {
			log.Debug(ctx, "A nested dirctory was encountered, skipping uploading the directory", "directory", info.Name())
			return nil
		}

		// Determine content type based on file extension
		contentType := getContentType(filepath.Ext(path))
		if contentType == "" {
			log.Debug(ctx, "Skipping uploading unsupported file type", "file", info.Name())
			return nil
		}

		// Open the file for reading
		file, err := os.Open(path)
		if err != nil {
			log.Error(ctx, "could not open the file in the directory", "file", file.Name(), "error", err)
			return err
		}
		defer file.Close()
		log.Debug(ctx, "opened file for uploading", "file", file.Name())
		// Upload the file
		objectName := baseFilename + "/" + info.Name()
		log.Debug(ctx, "printing upload params", "storage-bucket", cfg.String("minio.storage-bucket"), "object-name", objectName, "file-location", file.Name())
		_, err = minioClient.FPutObject(ctx, cfg.String("minio.storage-bucket"), objectName, file.Name(), minio.PutObjectOptions{
			ContentType: contentType,
		})
		if err != nil {
			log.Error(ctx, "Failed to upload file", "file", file.Name(), "error", err)
			return err
		}

		log.Debug(ctx, "Uploaded file into Minio/s3 storage", "object-name", objectName, "file-location", file.Name(), "content-type", contentType)
		return nil
	})

	if err != nil {
		log.Error(ctx, "Error walking the folder", "folder", currentDir, "error", err)
	}
	log.Info(ctx, "Folder uploaded successfully", "directory", currentDir)
}

func getContentType(extension string) string {
	switch extension {
	case ".mpd":
		return "application/dash+xml"
	case ".mp4":
		return "video/mp4"
	case ".m4s":
		return "video/iso.segment"
	default:
		return mime.TypeByExtension(extension) // Fallback to the standard MIME detection
	}
}
