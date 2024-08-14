package apihandler

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	domain "templify/pkg/domain/model"
	server "templify/pkg/server/generated"
	"templify/pkg/server/handler"

	"github.com/go-chi/render"
)

// Send a SMS with custom text
// (POST /sms)
func (ah *APIHandler) SendBasicSMS(w http.ResponseWriter, r *http.Request) {
	var smsRequest server.SMSSendRequest

	if err := handler.ReadRequestBody(w, r, &smsRequest); err != nil {
		return
	}

	err := ah.Usecase.SendSMS(smsRequest.ReceiverPhoneNumber, smsRequest.Message)
	if err != nil {
		handler.HandleError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	render.Status(r, http.StatusOK)
	render.PlainText(w, r, "SMS sent successfully")
}

func (ah *APIHandler) SendTemplatedSMS(w http.ResponseWriter, r *http.Request, templateName string) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		handler.HandleError(w, r, http.StatusBadRequest, "Reading Request Body failed")
		return
	}
	var smsRequest server.SMSTemplateSendRequest

	if err := json.Unmarshal(body, &smsRequest); err != nil {
		handler.HandleError(w, r, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	var filledTemplate string
	filledTemplate, err = ah.Usecase.GetFilledSMSTemplate(templateName, smsRequest.Placeholders)
	slog.With(
		"FilledTemplate", filledTemplate,
	).Debug("Filled template")
	if err != nil {
		handler.HandleError(w, r, http.StatusInternalServerError, "Error filling template")
		return
	}

	err = ah.Usecase.SendSMS(smsRequest.ReceiverPhoneNumber, filledTemplate)
	if err != nil {
		handler.HandleError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	render.Status(r, http.StatusOK)
	render.PlainText(w, r, "SMS sent successfully")
}

func (ah *APIHandler) AddNewSMSTemplate(w http.ResponseWriter, r *http.Request, templateName string) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		handler.HandleError(w, r, http.StatusBadRequest, "Reading Request Body failed")
		return
	}
	var SMSTempl server.SMSTemplate
	if err := json.Unmarshal(body, &SMSTempl); err != nil {
		handler.HandleError(w, r, http.StatusBadRequest, "Invalid JSON format")
		return
	}
	err = ah.Usecase.AddSMSTemplate(templateName, SMSTempl.TemplateString)
	if err != nil {
		handler.HandleError(w, r, http.StatusInternalServerError, fmt.Sprintf("Adding SMS template with name %v failed", templateName))
		return
	}
	resultString := fmt.Sprintf("Added SMS template with name %v", templateName)
	render.Status(r, http.StatusCreated)
	render.PlainText(w, r, resultString)
}

func (ah *APIHandler) FillSMSTemplate(w http.ResponseWriter, r *http.Request, templateName string) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Reading request body failed", http.StatusInternalServerError)
	}

	var templateFillRequest domain.SMSTemplateFillRequest

	if err := json.Unmarshal(body, &templateFillRequest); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	placeholders := templateFillRequest.Placeholders
	filledTemplate, err := ah.Usecase.GetFilledSMSTemplate(templateName, placeholders)
	if err != nil {
		handler.HandleError(w, r, http.StatusInternalServerError, "Error filling template")
		return
	}

	render.Status(r, http.StatusOK)
	render.PlainText(w, r, filledTemplate)
}

func (ah *APIHandler) GetSMSTemplateByName(w http.ResponseWriter, r *http.Request, templateName string) {
	templateDomain := &domain.Template{}
	var err error
	templateDomain, err = ah.Usecase.GetSMSTemplateByName(templateName)
	if err != nil {
		handler.HandleError(w, r, http.StatusInternalServerError, "Error getting template")
		return
	}
	if templateDomain.TemplateStr == "" {
		handler.HandleError(w, r, http.StatusNotFound, fmt.Sprintf("Template with name %s not found", templateName))
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, templateDomain)
}

func (ah *APIHandler) GetSMSTemplatePlaceholdersByName(w http.ResponseWriter, r *http.Request, templateName string) {
	templatePlaceholders, err := ah.Usecase.GetSMSPlaceholders(templateName)
	if err != nil {
		handler.HandleError(w, r, http.StatusInternalServerError, fmt.Sprintf("Getting placeholders for template %s failed", templateName))
		return
	}
	if len(templatePlaceholders) == 0 {
		handler.HandleError(w, r, http.StatusNotFound, fmt.Sprintf("No placeholders for template %s found", templateName))
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, templatePlaceholders)
}
