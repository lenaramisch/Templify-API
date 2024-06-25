package router

import (
	"example.SMSService.com/pkg/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func CreateRouter() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Post("/sms", handler.SMSPostRequest)
	router.Post("/email", handler.EmailPostRequest)
	return router
}
