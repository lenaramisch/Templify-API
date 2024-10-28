package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt"
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

func CheckIfAuthorised(w http.ResponseWriter, r *http.Request, requiredClaims map[string]any) bool {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		HandleError(w, r, http.StatusUnauthorized, "Authorization header missing")
		return false //unauthorized
	}
	tokenString := strings.TrimPrefix(authHeader, "Bearer")
	if tokenString == authHeader {
		HandleError(w, r, http.StatusUnauthorized, "Invalid token format, expected Bearer")
		return false //unauthorized
	}
	token, err := VerifyToken(tokenString)
	if err != nil || !token.Valid {
		HandleError(w, r, http.StatusUnauthorized, "Invalid token")
		return false //unauthorized
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		HandleError(w, r, http.StatusForbidden, "Access denied")
		return false //unauthorized
	}
	for key, value := range requiredClaims {
		if claims[key] != value {
			HandleError(w, r, http.StatusForbidden, "Access denied")
			return false //unauthorized
		}
	}
	return true //authorized
}

func VerifyToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return token, nil
	})
}
