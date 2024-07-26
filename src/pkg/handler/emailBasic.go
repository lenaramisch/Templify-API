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

func (ah *APIHandler) EmailPostRequest(res http.ResponseWriter, req *http.Request) {
	var emailRequestAPI EmailRequest

	decodeJSONBody(res, req, &emailRequestAPI)

	if emailRequestAPI.ToEmail == "" || emailRequestAPI.ToName == "" || emailRequestAPI.Subject == "" || emailRequestAPI.MessageBody == "" {
		http.Error(res, "Empty string content in either ToEmail, ToName, Subject or MessageBody", http.StatusBadRequest)
		return
	}
	emailRequestDomain := domain.EmailRequest{
		ToEmail:     emailRequestAPI.ToEmail,
		ToName:      emailRequestAPI.ToName,
		Subject:     emailRequestAPI.Subject,
		MessageBody: emailRequestAPI.MessageBody,
	}
	err := ah.usecase.SendRawEmail(&emailRequestDomain)
	if err != nil {
		handleError(res, req, http.StatusInternalServerError, err.Error())
		return
	}
	render.Status(req, http.StatusOK)
	render.PlainText(res, req, "Email sent successfully")
}

func (ah *APIHandler) TemplatePostRequest(res http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		handleError(res, req, http.StatusBadRequest, "Reading Request Body failed")
		return
	}
	MJMLString := string(body)

	templateName := chi.URLParam(req, "templateName")
	if templateName == "" {
		http.Error(res, "URL Param templateName empty", http.StatusBadRequest)
		return
	}

	err = ah.usecase.AddEmailTemplate(templateName, MJMLString)
	if err != nil {
		handleError(res, req, http.StatusInternalServerError, fmt.Sprintf("Adding template with name %v failed", templateName))
		return
	}
	resultString := fmt.Sprintf("Added template with name %v", templateName)
	render.Status(req, http.StatusCreated)
	render.PlainText(res, req, resultString)
}

func (ah *APIHandler) GetTemplatePlaceholdersRequest(res http.ResponseWriter, req *http.Request) {
	templateName := chi.URLParam(req, "templateName")
	if templateName == "" {
		http.Error(res, "URL Param templateName empty", http.StatusBadRequest)
		return
	}
	templatePlaceholders, err := ah.usecase.GetEmailPlaceholders(templateName)
	if err != nil {
		handleError(res, req, http.StatusInternalServerError, fmt.Sprintf("Getting placeholders for template %s failed", templateName))
		return
	}
	if len(templatePlaceholders) == 0 {
		handleError(res, req, http.StatusNotFound, fmt.Sprintf("No placeholders for template %s found", templateName))
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
	templateDomain, err = ah.usecase.GetEmailTemplateByName(templateName)
	if err != nil {
		handleError(res, req, http.StatusInternalServerError, "Error getting template")
		return
	}
	if templateDomain.TemplateStr == "" {
		handleError(res, req, http.StatusNotFound, fmt.Sprintf("Template with name %s not found", templateName))
		return
	}
	render.Status(req, http.StatusOK)
	render.JSON(res, req, templateDomain)

}

func (ah *APIHandler) PostTemplatePlaceholdersRequest(res http.ResponseWriter, req *http.Request) {
	templateName := chi.URLParam(req, "templateName")
	if templateName == "" {
		http.Error(res, "URL Param templateName empty", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, "Reading request body failed", http.StatusInternalServerError)
	}

	var templateFillRequest map[string]string

	if err := json.Unmarshal(body, &templateFillRequest); err != nil {
		http.Error(res, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	filledTemplate, err := ah.usecase.GetFilledMJMLTemplate(templateName, templateFillRequest)
	if err != nil {
		handleError(res, req, http.StatusInternalServerError, "Error filling template")
		return
	}

	render.Status(req, http.StatusOK)
	render.PlainText(res, req, filledTemplate)
}
