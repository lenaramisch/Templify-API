package apihandler

import (
	"templify/pkg/server/handler"

	"net/http"

	"github.com/go-chi/render"
)

// Download file from S3 bucket
// (GET /file/download/{fileName})
func (ah *APIHandler) GetDownloadFileURL(w http.ResponseWriter, r *http.Request, fileName string) {
	downloadURL, err := ah.Usecase.GetFileDownloadURL(fileName)
	if err != nil {
		handler.HandleInternalServerError(w, r, err, "Error getting download URL")
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, downloadURL)
}

// Upload file to S3 bucket for later use
// (POST /file/upload)
func (ah *APIHandler) GetUploadFileURL(w http.ResponseWriter, r *http.Request, fileName string) {
	uploadURL, err := ah.Usecase.GetFileUploadURL(fileName)
	if err != nil {
		handler.HandleInternalServerError(w, r, err, "Error getting upload URL")
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, uploadURL)
}
