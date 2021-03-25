package utils

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/defaults"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func UploadToS3(collector string, tempFile *os.File) error {

	cfg := defaults.Get().Config

	debug, _ := strconv.ParseBool(os.Getenv("DEBUG"))

	if debug {
		cfg.DisableSSL = aws.Bool(true)
		cfg.S3ForcePathStyle = aws.Bool(true)
	}

	if region, ok := os.LookupEnv("AWS_DEFAULT_REGION"); ok {
		cfg.Region = &region
	}

	if endpoint, ok := os.LookupEnv("AWS_DEFAULT_ENDPOINT"); ok {
		cfg.Endpoint = &endpoint
	}

	bucket := os.Getenv("AWS_DEFAULT_BUCKET")

	// The session the S3 Uploader will use
	sess := session.Must(session.NewSession(cfg))

	// Create an uploader with the session and default options
	uploader := s3manager.NewUploader(sess)

	key := path.Join(collector, filepath.Base(tempFile.Name()))

	upParams := &s3manager.UploadInput{
		Bucket: &bucket,
		Key:    &key,
		Body:   tempFile,
	}

	// Upload the file to S3.
	result, err := uploader.Upload(upParams)

	if err != nil {
		return fmt.Errorf("failed to upload file, %v", err)
	}

	fmt.Printf("file uploaded to, %s\n", result.Location)

	return nil
}
