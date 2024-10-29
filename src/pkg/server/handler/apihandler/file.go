package apihandler

import (
	"templify/pkg/server/handler"

	"net/http"

	"github.com/go-chi/render"
)

// Download file from S3 bucket
// (GET /file/download/{fileName})
func (ah *APIHandler) GetDownloadFileURL(w http.ResponseWriter, r *http.Request, fileName string) {
	requiredClaims := map[string]any{"role": "user"}
	checkedAuthorization := ah.Authorizer.CheckIfAuthorised(w, r, requiredClaims)
	if !checkedAuthorization {
		return
	}
	downloadURL, err := ah.Usecase.GetFileDownloadURL(fileName)
	if err != nil {
		handler.HandleErrors(w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, downloadURL)
}

// Upload file to S3 bucket for later use
// (POST /file/upload)
func (ah *APIHandler) GetUploadFileURL(w http.ResponseWriter, r *http.Request, fileName string) {
	requiredClaims := map[string]any{"role": "user"}
	checkedAuthorization := ah.Authorizer.CheckIfAuthorised(w, r, requiredClaims)
	if !checkedAuthorization {
		return
	}
	uploadURL, err := ah.Usecase.GetFileUploadURL(fileName)
	if err != nil {
		handler.HandleErrors(w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, uploadURL)
}
