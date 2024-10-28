package apihandler

import (
	"log/slog"
	domain "templify/pkg/domain/model"
	usecase "templify/pkg/domain/usecase"
	"templify/pkg/server/handler/authorisation"
)

type APIHandler struct {
	Usecase    *usecase.Usecase
	Authorizer *authorisation.Authorizer
	Info       *domain.Info
	log        *slog.Logger
	BaseURL    string
}

func NewAPIHandler(usecase *usecase.Usecase, authorizer *authorisation.Authorizer, info *domain.Info, logger *slog.Logger, baseURL string) *APIHandler {
	return &APIHandler{
		Usecase:    usecase,
		Authorizer: authorizer,
		Info:       info,
		log:        logger,
		BaseURL:    baseURL,
	}
}
