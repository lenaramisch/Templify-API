package router

import (
	"example.SMSService.com/pkg/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func CreateRouter(apiHandler *handler.APIHandler) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	// sms
	router.Post("/sms", apiHandler.SMSPostRequest)
	// email basic
	router.Post("/email", apiHandler.EmailPostRequest)

	// email templates
	router.Post("/email/templates/{templateName}", apiHandler.TemplatePostRequest)
	router.Get("/email/templates/{templateName}", apiHandler.GetTemplateByName)
	router.Get("/email/templates/{templateName}/placeholders", apiHandler.GetTemplatePlaceholdersRequest)

	// email mjml
	router.Post("/email/templates/{templateName}/placeholders", apiHandler.PostTemplatePlaceholdersRequest)
	router.Post("/email/templates/{templateName}/placeholders/attachments", apiHandler.SendMJMLWithAttachment)

	// email attachment
	router.Post("/email/attachments", apiHandler.SendEmailWithAttachment)

	// pdf
	router.Post("/pdf/templates/{templateName}", apiHandler.PDFTemplPostReq)
	router.Get("/pdf/templates/{templateName}", apiHandler.GetPDFTemplByName)
	router.Post("/pdf/templates/{templateName}/placeholders", apiHandler.GetFilledPDFTemplate)
	//TODO Implement route to get placeholders of PDF templ
	return router
}
