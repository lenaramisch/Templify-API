package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	domain "templify/pkg/domain/model"
	"time"

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

func HandleUnauthorizedError(w http.ResponseWriter, r *http.Request, err error, logMsg ...string) {
	if err == nil {
		err = errors.New("no error information supplied")
	}
	uniqueErrID := uuid.New().String()
	apiError := ErrorPayload{
		ErrorID:   uniqueErrID,
		Code:      401,
		Error:     err.Error(),
		ErrorType: "Unauthorized",
		Timestamp: time.Now().Format(time.RFC3339),
	}

	slog.With("error", err.Error()).With("logMessages", logMsg).Warn("Unauthorized Error")
	RespondWithJSON(w, r, http.StatusUnauthorized, apiError)
}

func HandleForbiddenError(w http.ResponseWriter, r *http.Request, err error, logMsg ...string) {
	if err == nil {
		err = errors.New("no error information supplied")
	}
	uniqueErrID := uuid.New().String()
	apiError := ErrorPayload{
		ErrorID:   uniqueErrID,
		Code:      403,
		Error:     err.Error(),
		ErrorType: "Forbidden",
		Timestamp: time.Now().Format(time.RFC3339),
	}

	slog.With("error", err.Error()).With("logMessages", logMsg).Warn("Forbidden Error")
	RespondWithJSON(w, r, http.StatusForbidden, apiError)
}

func HandleNotFoundError(w http.ResponseWriter, r *http.Request, err error, logMsg ...string) {
	if err == nil {
		err = errors.New("no error information supplied")
	}
	uniqueErrID := uuid.New().String()
	apiError := ErrorPayload{
		ErrorID:   uniqueErrID,
		Code:      404,
		Error:     err.Error(),
		ErrorType: "NotFound",
		Timestamp: time.Now().Format(time.RFC3339),
	}

	slog.With("error", err.Error()).With("logMessages", logMsg).Warn("Not Found Error")
	RespondWithJSON(w, r, http.StatusNotFound, apiError)
}

func HandleErrors(w http.ResponseWriter, r *http.Request, err error) {
	// HandleNotFoundError
	if errors.As(err, &domain.ErrorTemplateNotFound{}) ||
		errors.As(err, &domain.ErrorWorkflowNotFound{}) {
		HandleNotFoundError(w, r, err)
		return
	}
	// HandleBadRequestError
	if errors.As(err, &domain.ErrorPlaceholderMissing{}) ||
		errors.As(err, &domain.ErrorTemplateAlreadyExists{}) ||
		errors.As(err, &domain.ErrorGettingUploadURL{}) ||
		errors.As(err, &domain.ErrorWorkflowAlreadyExists{}) ||
		errors.As(err, &domain.ErrorAttachmentNameInvalid{}) {
		HandleBadRequestError(w, r, err)
		return
	}
	// HandleUnauthorizedError
	if errors.Is(err, domain.ErrAuthorizationHeaderMissing) ||
		errors.Is(err, domain.ErrInvalidTokenFormat) ||
		errors.Is(err, domain.ErrInvalidToken) {
		HandleUnauthorizedError(w, r, err)
		return
	}
	// HandleForbiddenError
	if errors.Is(err, domain.ErrAccessDenied) {
		HandleForbiddenError(w, r, err)
		return
	}
	// HandleInternalServerError
	HandleInternalServerError(w, r, err)
}
