package handler

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"unicode"

	"templify/pkg/domain"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type APIHandler struct {
	usecase *domain.Usecase
}

func NewAPIHandler(usecase *domain.Usecase) *APIHandler {
	return &APIHandler{
		usecase: usecase,
	}
}
