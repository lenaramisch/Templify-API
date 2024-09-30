package apihandler

import (
	"io"
	"strings"
	domain "templify/pkg/domain/model"
	"templify/pkg/server/handler"

	"net/http"

	"github.com/go-chi/render"
)

// Download file from S3 bucket
// (GET /file/download/{fileName})
func (ah *APIHandler) DownloadFile(w http.ResponseWriter, r *http.Request, fileName string) {
	splitString := strings.Split(fileName, ".")
	if len(splitString) != 2 {
		handler.HandleError(w, r, http.StatusBadRequest, "Invalid file name")
		return
	}
	fileDownloadRequest := domain.FileDownloadRequest{
		FileName:  splitString[0],
		Extension: splitString[1],
	}

	fileBytes, err := ah.Usecase.DownloadFile(fileDownloadRequest)
	if err != nil {
		handler.HandleError(w, r, http.StatusInternalServerError, "Download File failed")
		return
	}

	render.Status(r, http.StatusOK)
	render.Data(w, r, fileBytes)
}

// Upload file to S3 bucket for later use
// (POST /file/upload)
func (ah *APIHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// split fileName from extension
	splitString := strings.Split(handler.Filename, ".")

	// read file bytes
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Error reading the file", http.StatusBadRequest)
		return
	}

	domainFileUploadRequest := domain.FileUploadRequest{
		FileName:  splitString[0],
		Extension: splitString[1],
		FileBytes: fileBytes,
	}

	//call domain function
	err = ah.Usecase.UploadFile(domainFileUploadRequest)
	if err != nil {
		http.Error(w, "Error uploading the file", http.StatusInternalServerError)
		return
	}

	render.Status(r, http.StatusOK)
	render.PlainText(w, r, "File uploaded successfully")
}
