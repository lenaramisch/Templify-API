package router

import (
	"example.SMSService.com/pkg/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func CreateRouter(apiHandler *handler.APIHandler) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Post("/sms", apiHandler.SMSPostRequest)
	router.Post("/email", apiHandler.EmailPostRequest)
	return router
}
