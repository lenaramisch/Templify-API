package filemanager

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"net/url"
	domain "templify/pkg/domain/model"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	smithyendpoints "github.com/aws/smithy-go/endpoints"
)

type FileManagerConfig struct {
	BaseURL     string
	Port        string
	BucketName  string
	Region      string
	AccessKeyID string
	SecretKeyID string
}

type FileManager struct {
	config     FileManagerConfig
	log        *slog.Logger
	s3Client   *s3.Client
	s3psClient *s3.PresignClient
}

type minioResolver struct {
	BaseURL string
	Port    string
}

func (mr *minioResolver) ResolveEndpoint(ctx context.Context, params s3.EndpointParameters) (smithyendpoints.Endpoint, error) {
	u, err := url.Parse(mr.BaseURL + ":" + mr.Port)
	if err != nil {
		return smithyendpoints.Endpoint{}, err
	}
	return smithyendpoints.Endpoint{
		URI: *u,
	}, nil
}

// This filemanager uses s3 as a storage
func NewFileManagerService(fmCfg *FileManagerConfig, log *slog.Logger) *FileManager {
	staticProvider := credentials.NewStaticCredentialsProvider(
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
	}

	minioResolver := &minioResolver{
		BaseURL: fmCfg.BaseURL,
		Port:    fmCfg.Port,
	}

	client := s3.NewFromConfig(
		cfg,
		func(o *s3.Options) {
			o.UsePathStyle = true
			o.EndpointResolverV2 = minioResolver
		},
	)

	return &FileManager{
		config:     *fmCfg,
		log:        log,
		s3Client:   client,
		s3psClient: s3.NewPresignClient(client),
	}
}

func (fm *FileManager) GetFileDownloadURL(fileName string) (string, error) {
	request, err := fm.s3psClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
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

func (fm *FileManager) GetFileUploadURL(fileName string) (string, error) {
	request, err := fm.s3psClient.PresignPostObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(fm.config.BucketName),
		Key:    aws.String(fileName),
	}, func(options *s3.PresignPostOptions) {
		options.Expires = time.Duration(60*60*12) * time.Second
	})
	if err != nil {
		fm.log.Debug("Couldn't get a presigned post request to put %v:%v. Here's why: %v\n", fm.config.BucketName, fileName, slog.Any("error", err))
	}
	return request.URL, nil
}

// ListBuckets lists the buckets in the current account.
func (fm *FileManager) ListBuckets() ([]types.Bucket, error) {
	result, err := fm.s3Client.ListBuckets(context.Background(), &s3.ListBucketsInput{})
	if err != nil {
		fm.log.With("Error", err.Error()).Debug("Error listing buckets")
		return nil, err
	}
	return result.Buckets, err
}

// func (fm *FileManager) GetFileDownloadURL(fileName string) (string, error) {

// 	req, err := fm.s3psClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
// 		Bucket: aws.String(fm.config.BucketName),
// 		Key:    aws.String(fileName),
// 	}, s3.WithPresignExpires(time.Hour*12))

// 	if err != nil {
// 		fm.log.With("Error", err.Error()).Debug("Failed to sign request")
// 		return "", err
// 	}

// 	return req.URL, nil
// }

// func (fm *FileManager) GetFileUploadURL(fileName string) (string, error) {
// 	req, err := fm.s3psClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
// 		Bucket: aws.String(fm.config.BucketName),
// 		Key:    aws.String(fileName),
// 	}, s3.WithPresignExpires(time.Hour*12))

// 	if err != nil {
// 		fm.log.With("Error", err.Error()).Debug("Failed to sign request")
// 		return "", err
// 	}

// 	return req.URL, nil
// }

// ListFiles lists the files in a bucket.
func (fm *FileManager) ListFiles(bucketName string) ([]types.Object, error) {
	result, err := fm.s3Client.ListObjectsV2(context.Background(), &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		fm.log.With("Error", err.Error()).Debug("Error listing files")
		return nil, err
	}
	return result.Contents, err
}

// upload file to s3 using aws sdk
func (fm *FileManager) UploadFile(fileUploadRequest domain.FileUploadRequest) error {
	// create filereader for the file bytes
	fileReader := bytes.NewReader(fileUploadRequest.FileBytes)

	_, err := fm.s3Client.PutObject(context.Background(), &s3.PutObjectInput{
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

func (fm *FileManager) DownloadFile(fileDownloadRequest domain.FileDownloadRequest) ([]byte, error) {
	objectKey := fileDownloadRequest.FileName + "." + fileDownloadRequest.Extension
	bucketName := fm.config.BucketName
	if fileDownloadRequest.BucketName != nil {
		bucketName = *fileDownloadRequest.BucketName
	}
	// check if file exists
	_, err := fm.s3Client.HeadObject(context.Background(), &s3.HeadObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		fm.log.With("Error", err.Error()).Debug("Error checking if file exists")
		return nil, err
	}
	result, err := fm.s3Client.GetObject(context.Background(), &s3.GetObjectInput{
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
