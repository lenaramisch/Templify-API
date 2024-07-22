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
	router.Post("/templates/{templateName}", apiHandler.TemplatePostRequest)
	router.Get("/templates/{templateName}", apiHandler.GetTemplateByName)
	router.Get("/templates/{templateName}/placeholders", apiHandler.GetTemplatePlaceholdersRequest)
	router.Post("/templates/{templateName}/placeholders", apiHandler.PostTemplatePlaceholdersRequest)
	router.Post("/email/attachments", apiHandler.EmailPostRequestAttm)
	router.Post("/templates/{templateName}/placeholders/attachments", apiHandler.PostTmplPlaceholdersAttm)
	router.Post("/templates/pdf/{templateName}", apiHandler.PDFTemplPostReq)
	router.Get("/templates/pdf/{templateName}", apiHandler.GetPDFTemplByName)
	router.Post("/templates/pdf/{templateName}/placeholders", apiHandler.PostPDFTemplPlaceholdersRequest)
	//TODO Implement route to get placeholders of PDF templ
	return router
}
