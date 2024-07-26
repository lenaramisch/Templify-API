package apihandler

import (
	"log/slog"
	domain "templify/pkg/domain/model"
	usecase "templify/pkg/domain/usecase"
)

type APIHandler struct {
	Usecase *usecase.Usecase
	Info    *domain.Info
	log     *slog.Logger
	BaseURL string
}

func NewAPIHandler(usecase *usecase.Usecase, info *domain.Info, logger *slog.Logger, baseURL string) *APIHandler {
	return &APIHandler{
		Usecase: usecase,
		Info:    info,
		log:     logger,
		BaseURL: baseURL,
	}
}
