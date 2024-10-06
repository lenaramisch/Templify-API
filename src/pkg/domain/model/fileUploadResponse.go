package domain

type FileUploadResponse struct {
	UploadURL string
	Values    *map[string]string
}
