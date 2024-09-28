package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"
	"unicode"

	"github.com/go-chi/render"
	"github.com/google/uuid"
)

// ErrorPayload defines the structure of the error response payload.
type ErrorPayload struct {
	ErrorID   string `json:"errorId"`
	Code      int    `json:"code"`
	Error     string `json:"error"`
	ErrorType string `json:"errorType"`
	Timestamp string `json:"timestamp"`
}

// HandleInternalServerError is a convenient method to log and handle internal server errors.
func HandleInternalServerError(w http.ResponseWriter, r *http.Request, err error, logMsg ...string) {
	if err == nil {
		err = errors.New("no error information supplied")
	}
	uniqueErrID := uuid.New().String()
	apiError := ErrorPayload{
		ErrorID:   uniqueErrID,
		Code:      500,
		Error:     err.Error(),
		ErrorType: "InternalServerError", // Assuming this is the type string you want
		Timestamp: time.Now().Format(time.RFC3339),
	}

	slog.With("error", err.Error()).With("logMessages", logMsg).Error("Internal Server Error")
	RespondWithJSON(w, r, http.StatusInternalServerError, apiError)
}

// HandleBadRequestError is a convenient method to log and handle bad request errors.
func HandleBadRequestError(w http.ResponseWriter, r *http.Request, err error, logMsg ...string) {
	if err == nil {
		err = errors.New("no error information supplied")
	}
	uniqueErrID := uuid.New().String()
	apiError := ErrorPayload{
		ErrorID:   uniqueErrID,
		Code:      400,
		Error:     err.Error(),
		ErrorType: "BadRequest", // Assuming this is the type string you want
		Timestamp: time.Now().Format(time.RFC3339),
	}

	slog.With("error", err.Error()).With("logMessages", logMsg).Warn("Bad Request Error")
	RespondWithJSON(w, r, http.StatusBadRequest, apiError)
}

func RespondWithJSON(w http.ResponseWriter, r *http.Request, status int, v any) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		HandleInternalServerError(w, r, err, "Failed to encode response")
	}
}

func ReadRequestBody(w http.ResponseWriter, r *http.Request, v any) error {
	err := json.NewDecoder(r.Body).Decode(v)
	if err != nil {
		HandleBadRequestError(w, r, err, "Failed to parse request body")
	}
	return err
}

func HandleError(w http.ResponseWriter, r *http.Request, statusCode int, message string) {
	render.Status(r, statusCode)
	render.PlainText(w, r, message)
}

func FormToCapitalPlaceholders(r *http.Request) {
	form := r.MultipartForm
	placeholders := map[string]string{}
	for key, values := range form.Value {
		if len(key) > 0 && unicode.IsUpper(rune(key[0])) {
			placeholders[key] = values[0]
		}
	}
}
