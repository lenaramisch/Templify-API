package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"unicode"

	"example.SMSService.com/pkg/domain"
	"github.com/go-chi/render"
)

func handleError(res http.ResponseWriter, req *http.Request, statusCode int, message string) {
	render.Status(req, statusCode)
	render.PlainText(res, req, message)
}

func decodeJSONBody(res http.ResponseWriter, req *http.Request, dst any) error {
	err := json.NewDecoder(req.Body).Decode(dst)
	if err != nil {
		handleError(res, req, http.StatusBadRequest, "Invalid JSON: "+err.Error())
		return err
	}
	return nil
}

func readMultipartFileAsBytes(req *http.Request, res http.ResponseWriter) (*domain.AttachmentInfo, error) {
	// Parse the multipart form, with a maximum memory of 32 MB for storing file parts in memory
	err := req.ParseMultipartForm(32 << 20) // 32MB
	if err != nil {
		return nil, err
	}

	file, handler, err := req.FormFile("file")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		handleError(res, req, http.StatusInternalServerError, "Error reading file content")
		return nil, err
	}

	// Get file type string
	lastDotIndex := strings.LastIndex(handler.Filename, ".")
	if lastDotIndex == -1 {
		http.Error(res, "File name has to end with file extension, i.e. '.txt'", http.StatusBadRequest)
		return nil, err
	}

	fileTypeString := handler.Filename[lastDotIndex+1:]
	fileName := handler.Filename

	return &domain.AttachmentInfo{
		FileName:      fileName,
		FileExtension: fileTypeString,
		Content:       fileBytes,
	}, nil
}

func formToCapitalPlaceholders(req *http.Request) {
	form := req.MultipartForm
	placeholders := map[string]string{}
	for key, values := range form.Value {
		if len(key) > 0 && unicode.IsUpper(rune(key[0])) {
			placeholders[key] = values[0]
		}
	}
}
