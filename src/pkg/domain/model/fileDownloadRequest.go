package domain

type FileDownloadRequest struct {
	FileName   string
	Extension  string
	BucketName *string
}
