package filemanager

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	domain "templify/pkg/domain/model"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type FileManager struct {
	config   Config
	log      *slog.Logger
	s3Client *s3.Client
}

// This filemanager uses s3 as a storage
func NewFileManager(config *Config, log *slog.Logger) *FileManager {
	staticProvider := credentials.NewStaticCredentialsProvider(
		config.AccessKeyID,
		config.SecretKeyID,
		"",
	)

	awsCfg := aws.Config{
		Region:       config.Region,
		BaseEndpoint: aws.String(config.BaseURL + ":" + config.Port),
		Credentials:  staticProvider,
	}
	client := s3.NewFromConfig(awsCfg, nil)

	return &FileManager{
		config:   *config,
		log:      log,
		s3Client: client,
	}
}

// ListBuckets lists the buckets in the current account.
func (fm *FileManager) ListBuckets() ([]types.Bucket, error) {
	result, err := fm.s3Client.ListBuckets(context.Background(), &s3.ListBucketsInput{})
	if err != nil {
		fm.log.With("Error", err.Error()).Debug("Error listing buckets")
	}
	return result.Buckets, err
}

// upload file to s3 using aws sdk
func (fm *FileManager) UploadFile(fileUploadRequest domain.FileUploadRequest) error {
	// create filereader for the file bytes
	fileReader := bytes.NewReader(fileUploadRequest.File)

	_, err := fm.s3Client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String(fm.config.BucketName),
		Key:    aws.String(fileUploadRequest.FileName + "." + fileUploadRequest.Extension),
		Body:   fileReader,
	})
	if err != nil {
		fm.log.With("Error", err.Error()).Debug("Error uploading file")
	}
	return nil
}

func (fm *FileManager) DownloadFile(fileDownloadRequest domain.FileDownloadRequest) ([]byte, error) {
	objectKey := fileDownloadRequest.FileName + "." + fileDownloadRequest.Extension
	result, err := fm.s3Client.GetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String(fm.config.BucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		fm.log.With("Error", err.Error()).Debug("Error downloading file")
		return nil, err
	}
	downloadedFile := make([]byte, 0)
	_, err = result.Body.Read(downloadedFile)
	if err != nil && err != io.EOF {
		fm.log.With("Error", err.Error()).Debug("Error reading file")
		return nil, err
	}
	defer result.Body.Close()
	return downloadedFile, nil
}
