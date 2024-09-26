package domain

type FileUploadRequest struct {
	FileName  string
	Extension string
	File      []byte
}
