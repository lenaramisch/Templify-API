package handler

import (
	"encoding/json"
	"fmt"
	"io"
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

func (ah *APIHandler) TemplatePostRequest(res http.ResponseWriter, req *http.Request) {

	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, "Reading Request Body failed", http.StatusBadRequest)
		return
	}
	MJMLString := string(body)

	templateName := chi.URLParam(req, "templateName")
	if templateName == "" {
		http.Error(res, "URL Param templateName empty", http.StatusBadRequest)
		return
	}

	err = ah.usecase.AddTemplate(templateName, MJMLString)
	if err != nil {
		http.Error(res, fmt.Sprintf("Adding template with name %v failed", templateName), http.StatusInternalServerError)
		return
	}
	resultString := fmt.Sprintf("Added template with name %v", templateName)
	render.Status(req, http.StatusOK)
	render.PlainText(res, req, resultString)
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

func (ah *APIHandler) GetTemplateByName(res http.ResponseWriter, req *http.Request) {
	templateName := chi.URLParam(req, "templateName")
	if templateName == "" {
		http.Error(res, "URL Param templateName empty", http.StatusBadRequest)
		return
	}
	templateDomain := &domain.Template{}
	var err error
	templateDomain, err = ah.usecase.GetTemplateByName(templateName)
	if err != nil {
		http.Error(res, "Error getting template", http.StatusInternalServerError)
		return
	}
	render.Status(req, http.StatusOK)
	render.JSON(res, req, templateDomain)

}

func (ah *APIHandler) PostTemplatePlacehholdersRequest(res http.ResponseWriter, req *http.Request) {
	templateName := chi.URLParam(req, "templateName")
	if templateName == "" {
		http.Error(res, "URL Param templateName empty", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, "Reading request body failed", http.StatusInternalServerError)
	}

	values := map[string]any{}

	if err := json.Unmarshal(body, &values); err != nil {
		http.Error(res, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	filledTemplate, err := ah.usecase.FillTemplatePlaceholders(templateName, values)
	if err != nil {
		http.Error(res, "Error filling template", http.StatusInternalServerError)
		return
	}

	render.Status(req, http.StatusOK)
	render.PlainText(res, req, filledTemplate)
}
