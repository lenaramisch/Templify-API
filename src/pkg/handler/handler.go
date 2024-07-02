package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"example.SMSService.com/pkg/domain"
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

func (ah *APIHandler) SMSPostRequest(res http.ResponseWriter, req *http.Request) {
	var smsRequest SmsRequest

	err := json.NewDecoder(req.Body).Decode(&smsRequest)
	if err != nil {
		render.Status(req, http.StatusBadRequest)
		render.PlainText(res, req, "You messed up: "+err.Error())
		return
	}

	if smsRequest.ToNumber == "" || smsRequest.MessageBody == "" {
		http.Error(res, "Empty string content in either ToNumber or MessageBody", http.StatusBadRequest)
		return
	}
	err = ah.usecase.SendSMS(smsRequest.ToNumber, smsRequest.MessageBody)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	render.Status(req, http.StatusOK)
	render.PlainText(res, req, "SMS sent successfully")
}

func (ah *APIHandler) EmailPostRequest(res http.ResponseWriter, req *http.Request) {
	var emailRequest EmailRequest

	err := json.NewDecoder(req.Body).Decode(&emailRequest)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if emailRequest.ToEmail == "" || emailRequest.ToName == "" || emailRequest.Subject == "" || emailRequest.MessageBody == "" {
		http.Error(res, "Empty string content in either ToEmail, ToName, Subject or MessageBody", http.StatusBadRequest)
		return
	}
	err = ah.usecase.SendEmail(emailRequest.ToEmail, emailRequest.ToName, emailRequest.Subject, emailRequest.MessageBody)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	render.Status(req, http.StatusOK)
	render.PlainText(res, req, "Email sent successfully")
}

// TODO add MJML functions
func (ah *APIHandler) TemplatePostRequest(res http.ResponseWriter, req *http.Request) {
	return
}
func (ah *APIHandler) GetTemplatePlaceholdersRequest(res http.ResponseWriter, req *http.Request) {
	templateName := chi.URLParam(req, "templateName")
	if templateName == "" {
		http.Error(res, "URL Param templateName empty", http.StatusBadRequest)
		return
	}
	templatePlaceholders, err := ah.usecase.GetTemplatePlaceholders(templateName)
	if err != nil {
		http.Error(res, fmt.Sprintf("Getting placeholders for template %v failed", templateName), http.StatusInternalServerError)
		return
	}
	render.Status(req, http.StatusOK)
	render.JSON(res, req, templatePlaceholders)
}

func (ah *APIHandler) PostTemplatePlacehholdersRequest(res http.ResponseWriter, req *http.Request) {
	//get templateName from Request Params
	templateName := chi.URLParam(req, "templateName")
	if templateName == "" {
		http.Error(res, "URL Param templateName empty", http.StatusBadRequest)
		return
	}

	//TODO read values for placeholders (req body)
	values := map[string]any{}

	filledTemplate, err := ah.usecase.FillTemplatePlaceholders(templateName, values)
	if err != nil {
		http.Error(res, "Error filling template", http.StatusInternalServerError)
		return
	}
	render.Status(req, http.StatusOK)
	render.PlainText(res, req, filledTemplate)
}
