package filemanager

import (
	"bytes"
	"context"
	"log/slog"
	domain "templify/pkg/domain/model"
	"time"

	minio "github.com/minio/minio-go/v7"
	minioCreds "github.com/minio/minio-go/v7/pkg/credentials"
)

type FileManagerMinio struct {
	config      FileManagerConfig
	log         *slog.Logger
	minioClient *minio.Client
}

// This filemanager uses s3 as a storage
func NewFileManagerMinioService(fmCfg *FileManagerConfig, log *slog.Logger) *FileManagerMinio {
	endpoint := fmCfg.BaseURL + ":" + fmCfg.Port
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  minioCreds.NewStaticV4(fmCfg.AccessKeyID, fmCfg.SecretKeyID, ""),
		Secure: false,
	})
	if err != nil {
		log.With("Error", err.Error()).Debug("Error creating minio client")
	}
	return &FileManagerMinio{
		config:      *fmCfg,
		log:         log,
		minioClient: minioClient,
	}
}

func (fm *FileManagerMinio) GetFileDownloadURL(fileName string) (string, error) {
	// get a presigned download URL for the file
	presignedURL, err := fm.minioClient.PresignedGetObject(context.Background(), fm.config.BucketName, fileName, time.Hour, nil)
	if err != nil {
		fm.log.With("Error", err.Error()).Debug("Error getting presigned url")
		return "", err
	}
	return presignedURL.String(), nil
}

func (fm *FileManagerMinio) GetFileUploadURL(fileName string) (*domain.FileUploadResponse, error) {
	// get a presigned upload URL for the file
	presignedURL, err := fm.minioClient.PresignedPutObject(context.Background(), fm.config.BucketName, fileName, time.Hour)
	if err != nil {
		fm.log.With("Error", err.Error()).Debug("Error getting presigned url")
		return nil, err
	}
	return &domain.FileUploadResponse{
		UploadURL: presignedURL.String(),
	}, nil
}

func (fm *FileManagerMinio) ListBuckets() ([]string, error) {
	bucketInfo, err := fm.minioClient.ListBuckets(context.Background())
	if err != nil {
		fm.log.With("Error", err.Error()).Debug("Error listing buckets")
		return nil, err
	}
	buckets := make([]string, 0)
	for _, bucket := range bucketInfo {
		buckets = append(buckets, bucket.Name)
	}
	return buckets, nil
}

func (fm *FileManagerMinio) ListFiles(bucketName string) ([]string, error) {
	fileInfo := fm.minioClient.ListObjects(context.Background(), bucketName, minio.ListObjectsOptions{})
	objects := make([]string, 0)
	for object := range fileInfo {
		objects = append(objects, object.Key)
	}
	return objects, nil
}

func (fm *FileManagerMinio) UploadFile(fileUploadRequest domain.FileUploadRequest) error {
	// create reader for fileBytes
	reader := bytes.NewReader(fileUploadRequest.FileBytes)
	_, err := fm.minioClient.PutObject(context.Background(), fm.config.BucketName, fileUploadRequest.FileName+"."+fileUploadRequest.Extension, reader, -1, minio.PutObjectOptions{})
	if err != nil {
		fm.log.With("Error", err.Error()).Debug("Error uploading file")
		return err
	}
	return nil
}

func (fm *FileManagerMinio) DownloadFile(fileDownloadRequest domain.FileDownloadRequest) ([]byte, error) {
	object, err := fm.minioClient.GetObject(context.Background(), fm.config.BucketName, fileDownloadRequest.FileName+"."+fileDownloadRequest.Extension, minio.GetObjectOptions{})
	if err != nil {
		fm.log.With("Error", err.Error()).Debug("Error downloading file")
		return nil, err
	}
	defer object.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(object)
	return buf.Bytes(), nil
}
