package utils

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func UploadToS3(collector string, tempFile *os.File) error {

	bucket := os.Getenv("AWS_DEFAULT_BUCKET")
	key := path.Join(collector, filepath.Base(tempFile.Name()))

	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		return err
	}

	// Create an Amazon S3 service client
	client := s3.NewFromConfig(cfg)

	// Upload the file to S3.
	uploader := manager.NewUploader(client)

	result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   tempFile,
	})

	if err != nil {
		return fmt.Errorf("failed to upload file, %v", err)
	}

	fmt.Printf("file uploaded to %s", result.Location)

	return nil
}
