package filemanager

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	domain "templify/pkg/domain/model"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	awsCreds "github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type FileManagerAWS struct {
	config        FileManagerConfig
	log           *slog.Logger
	awsS3Client   *s3.Client
	awsS3PsClient *s3.PresignClient
}

// This filemanager uses s3 as a storage
func NewFileManagerAWSService(fmCfg *FileManagerConfig, log *slog.Logger) *FileManagerAWS {
	staticProvider := awsCreds.NewStaticCredentialsProvider(
		fmCfg.AccessKeyID,
		fmCfg.SecretKeyID,
		"",
	)

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(fmCfg.Region),
		config.WithCredentialsProvider(staticProvider),
	)
	if err != nil {
		log.With("Error", err.Error()).Debug("Error loading default config")
		return nil
	}

	client := s3.NewFromConfig(
		cfg,
		func(o *s3.Options) {
			o.UsePathStyle = true
		},
	)

	return &FileManagerAWS{
		config:        *fmCfg,
		log:           log,
		awsS3Client:   client,
		awsS3PsClient: s3.NewPresignClient(client),
	}
}

func (fm *FileManagerAWS) GetFileDownloadURL(fileName string) (string, error) {
	request, err := fm.awsS3PsClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(fm.config.BucketName),
		Key:    aws.String(fileName),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(60 * 60 * 12 * (time.Second))
	})
	if err != nil {
		fm.log.Debug("Couldn't get a presigned request to get %v:%v. Here's why: %v\n",
			fm.config.BucketName, fileName, slog.Any("error", err))
	}
	return request.URL, err
}

func (fm *FileManagerAWS) GetFileUploadURL(fileName string) (*domain.FileUploadResponse, error) {
	request, err := fm.awsS3PsClient.PresignPostObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(fm.config.BucketName),
		Key:    aws.String(fileName),
	}, func(options *s3.PresignPostOptions) {
		options.Expires = time.Duration(60*60*12) * time.Second
	})
	if err != nil {
		fm.log.Debug("Couldn't get a presigned post request to put %v:%v. Here's why: %v\n", fm.config.BucketName, fileName, slog.Any("error", err))
	}

	return &domain.FileUploadResponse{
		UploadURL: request.URL,
		Values:    &request.Values,
	}, nil
}

// ListBuckets lists the buckets in the current account.
func (fm *FileManagerAWS) ListBuckets() ([]string, error) {
	result, err := fm.awsS3Client.ListBuckets(context.Background(), &s3.ListBucketsInput{})
	if err != nil {
		fm.log.With("Error", err.Error()).Debug("Error listing buckets")
		return nil, err
	}
	buckets := make([]string, 0)
	for _, bucket := range result.Buckets {
		buckets = append(buckets, *bucket.Name)
	}
	return buckets, nil
}

// ListFiles lists the files in a bucket.
func (fm *FileManagerAWS) ListFiles(bucketName string) ([]string, error) {
	result, err := fm.awsS3Client.ListObjectsV2(context.Background(), &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		fm.log.With("Error", err.Error()).Debug("Error listing files")
		return nil, err
	}
	files := make([]string, 0)
	for _, file := range result.Contents {
		files = append(files, *file.Key)
	}
	return files, nil
}

// upload file to s3 using aws sdk
func (fm *FileManagerAWS) UploadFile(fileUploadRequest domain.FileUploadRequest) error {
	// create filereader for the file bytes
	fileReader := bytes.NewReader(fileUploadRequest.FileBytes)

	_, err := fm.awsS3Client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String(fm.config.BucketName),
		Key:    aws.String(fileUploadRequest.FileName + "." + fileUploadRequest.Extension),
		Body:   fileReader,
	})
	if err != nil {
		fm.log.With("Error", err.Error()).Debug("Error uploading file")
		return err
	}
	return nil
}

func (fm *FileManagerAWS) DownloadFile(fileDownloadRequest domain.FileDownloadRequest) ([]byte, error) {
	objectKey := fileDownloadRequest.FileName + "." + fileDownloadRequest.Extension
	bucketName := fm.config.BucketName
	if fileDownloadRequest.BucketName != nil {
		bucketName = *fileDownloadRequest.BucketName
	}
	// check if file exists
	_, err := fm.awsS3Client.HeadObject(context.Background(), &s3.HeadObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		fm.log.With("Error", err.Error()).Debug("Error checking if file exists")
		return nil, err
	}
	result, err := fm.awsS3Client.GetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	defer result.Body.Close()
	if err != nil {
		fm.log.With("Error", err.Error()).Debug("Error downloading file")
		return nil, err
	}
	downloadedFile, err := io.ReadAll(result.Body)
	if err != nil {
		fm.log.With("Error", err.Error()).Debug("Error reading file")
		return nil, err
	}
	return downloadedFile, nil
}
