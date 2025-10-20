package cloudstorage

import (
	"context"
	"fmt"
	"mime"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	// "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/knadh/koanf/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"hangout.com/core/storage-service/files"
	"hangout.com/core/storage-service/logger"
)

func Connect(workerId int, ctx context.Context, cfg *koanf.Koanf, log logger.Log) (*s3.Client, error) {
	log.Info(ctx, "connecting to AWS S3", "worker-id", workerId)

	region := cfg.String("aws.region")

	// Load default AWS config (automatically uses instance IAM role)
	awsCfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		log.Error(ctx, "could not load AWS configuration", "error", err)
		return nil, err
	}

	s3Client := s3.NewFromConfig(awsCfg)

	uploadBucket := cfg.String("s3.upload-bucket")
	storageBucket := cfg.String("s3.storage-bucket")

	checkAndCreateBucket(ctx, s3Client, uploadBucket, log, workerId)
	checkAndCreateBucket(ctx, s3Client, storageBucket, log, workerId)

	return s3Client, nil
}

// helper to ensure a bucket exists
func checkAndCreateBucket(ctx context.Context, client *s3.Client, bucket string, log logger.Log, workerId int) {
	_, err := client.HeadBucket(ctx, &s3.HeadBucketInput{Bucket: aws.String(bucket)})
	if err == nil {
		log.Debug(ctx, "Bucket exists", "bucket", bucket, "worker-id", workerId)
		return
	}

	var notFound bool
	// var ae *types.NotFound
	if strings.Contains(err.Error(), "NotFound") {
		notFound = true
	}

	if notFound {
		log.Info(ctx, "Bucket does not exist, creating", "bucket", bucket, "worker-id", workerId)
		_, err = client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket: aws.String(bucket),
		})
		if err != nil {
			log.Error(ctx, "Failed to create bucket", "bucket", bucket, "error", err, "worker-id", workerId)
		} else {
			log.Info(ctx, "Bucket created successfully", "bucket", bucket, "worker-id", workerId)
		}
	} else {
		log.Error(ctx, "Error checking bucket", "bucket", bucket, "error", err, "worker-id", workerId)
	}
}
func Download(ctx context.Context, s3Client *s3.Client, file *files.File, cfg *koanf.Koanf, log logger.Log) {
	tr := otel.Tracer("hangout.storage.cloudstorage")
	ctx, span := tr.Start(ctx, "DownloadFile")
	defer span.End()

	span.SetAttributes(
		attribute.String("file.name", file.Filename),
		attribute.Int("file.userId", int(file.UserId)),
	)

	log.Info(ctx, "Downloading file from S3", "file", file.Filename)

	downloadBucket := cfg.String("s3.upload-bucket")
	outputPath := "/tmp/" + file.Filename

	out, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(downloadBucket),
		Key:    aws.String(file.Filename),
	})
	if err != nil {
		log.Error(ctx, "Error while downloading file", "file", file.Filename, "error", err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return
	}
	defer out.Body.Close()

	f, err := os.Create(outputPath)
	if err != nil {
		log.Error(ctx, "Could not create local file", "path", outputPath, "error", err)
		span.RecordError(err)
		return
	}
	defer f.Close()

	_, err = f.ReadFrom(out.Body)
	if err != nil {
		log.Error(ctx, "Error saving file", "file", outputPath, "error", err)
		span.RecordError(err)
		return
	}

	log.Info(ctx, "File downloaded successfully", "path", outputPath)
}

func UploadDir(ctx context.Context, s3Client *s3.Client, event *files.File, cfg *koanf.Koanf, log logger.Log) {
	tr := otel.Tracer("hangout.storage.cloudstorage")
	ctx, span := tr.Start(ctx, "UploadDir")
	defer span.End()

	span.SetAttributes(
		attribute.String("file.name", event.Filename),
		attribute.Int("file.userId", int(event.UserId)),
	)

	baseFilename := strings.Split(event.Filename, ".")[0]
	currentDir := "/tmp/" + baseFilename
	storageBucket := cfg.String("s3.storage-bucket")

	log.Info(ctx, "Starting to upload directory to S3", "directory", currentDir)

	err := filepath.Walk(currentDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Error(ctx, "Error traversing file", "file", path, "error", err)
			return err
		}

		if info.IsDir() {
			return nil
		}

		contentType := getContentType(filepath.Ext(path))
		if contentType == "" {
			log.Debug(ctx, "Skipping unsupported file type", "file", info.Name())
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			log.Error(ctx, "Could not open file", "file", path, "error", err)
			return err
		}
		defer f.Close()

		objectKey := fmt.Sprintf("%s/%s", baseFilename, info.Name())

		log.Debug(ctx, "Uploading file", "bucket", storageBucket, "key", objectKey)

		_, err = s3Client.PutObject(ctx, &s3.PutObjectInput{
			Bucket:      aws.String(storageBucket),
			Key:         aws.String(objectKey),
			Body:        f,
			ContentType: aws.String(contentType),
		})
		if err != nil {
			log.Error(ctx, "Failed to upload file", "file", path, "error", err)
			return err
		}

		log.Debug(ctx, "Uploaded file", "key", objectKey)
		return nil
	})

	if err != nil {
		log.Error(ctx, "Error walking the folder", "folder", currentDir, "error", err)
	} else {
		log.Info(ctx, "Folder uploaded successfully", "directory", currentDir)
	}
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
