package apihandler

import (
	"fmt"
	"net/http"
	domain "templify/pkg/domain/model"
	server "templify/pkg/server/generated"
	"templify/pkg/server/handler"

	"github.com/go-chi/render"
)

// Send a SMS with custom text
// (POST /sms/basic/send)
func (ah *APIHandler) SendBasicSMS(w http.ResponseWriter, r *http.Request) {
	requiredClaims := map[string]any{"role": "user"}
	checkedAuthorization := ah.Authorizer.CheckIfAuthorised(w, r, requiredClaims)
	if !checkedAuthorization {
		return
	}

	var smsRequest server.SMSSendRequest
	if err := handler.ReadRequestBody(w, r, &smsRequest); err != nil {
		return
	}

	domainSmsRequest := domain.SmsRequest{
		ToNumber:    smsRequest.ReceiverPhoneNumber,
		MessageBody: smsRequest.Message,
	}

	err := ah.Usecase.SendSMS(domainSmsRequest)
	if err != nil {
		handler.HandleErrors(w, r, err)
		return
	}
	render.Status(r, http.StatusOK)
	render.PlainText(w, r, "SMS sent successfully")
}

// Send a SMS with template
// (POST /sms/templates/{templateName}/send)
func (ah *APIHandler) SendTemplatedSMS(w http.ResponseWriter, r *http.Request, templateName string) {
	requiredClaims := map[string]any{"role": "user"}
	checkedAuthorization := ah.Authorizer.CheckIfAuthorised(w, r, requiredClaims)
	if !checkedAuthorization {
		return
	}

	var smsRequest server.SMSTemplateSendRequest
	if err := handler.ReadRequestBody(w, r, &smsRequest); err != nil {
		return
	}

	smsFillRequest := domain.SMSTemplateFillRequest{
		ReceiverPhoneNumber: smsRequest.ReceiverPhoneNumber,
		Placeholders:        smsRequest.Placeholders,
	}

	err := ah.Usecase.SendTemplatedSMS(r.Context(), smsFillRequest, templateName)
	if err != nil {
		handler.HandleErrors(w, r, err)
		return
	}
	render.Status(r, http.StatusOK)
	render.PlainText(w, r, "SMS sent successfully")
}

// Add a new SMS template
// (POST /sms/templates/{templateName})
func (ah *APIHandler) AddNewSMSTemplate(w http.ResponseWriter, r *http.Request, templateName string) {
	requiredClaims := map[string]any{"role": "user"}
	checkedAuthorization := ah.Authorizer.CheckIfAuthorised(w, r, requiredClaims)
	if !checkedAuthorization {
		return
	}
	var SMSTempl server.SMSTemplate
	err := handler.ReadRequestBody(w, r, &SMSTempl)
	if err != nil {
		return
	}

	templateDomain := &domain.Template{
		Name:        templateName,
		TemplateStr: SMSTempl.TemplateString,
	}

	err = ah.Usecase.AddSMSTemplate(r.Context(), templateDomain)
	if err != nil {
		handler.HandleErrors(w, r, err)
		return
	}
	resultString := fmt.Sprintf("Added SMS template with name %v", templateName)
	render.Status(r, http.StatusCreated)
	render.PlainText(w, r, resultString)
}

// Fill placeholders of SMS template
// (POST /sms/templates/{templateName}/fill)
func (ah *APIHandler) FillSMSTemplate(w http.ResponseWriter, r *http.Request, templateName string) {
	requiredClaims := map[string]any{"role": "user"}
	checkedAuthorization := ah.Authorizer.CheckIfAuthorised(w, r, requiredClaims)
	if !checkedAuthorization {
		return
	}

	var templateFillRequest domain.SMSTemplateFillRequest
	err := handler.ReadRequestBody(w, r, &templateFillRequest)
	if err != nil {
		return
	}

	placeholders := templateFillRequest.Placeholders
	filledTemplate, err := ah.Usecase.GetFilledSMSTemplate(r.Context(), templateName, placeholders)
	if err != nil {
		handler.HandleErrors(w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
	render.PlainText(w, r, filledTemplate)
}

// Get SMS template by name
// (GET /sms/templates/{templateName})
func (ah *APIHandler) GetSMSTemplateByName(w http.ResponseWriter, r *http.Request, templateName string) {
	requiredClaims := map[string]any{"role": "user"}
	checkedAuthorization := ah.Authorizer.CheckIfAuthorised(w, r, requiredClaims)
	if !checkedAuthorization {
		return
	}
	templateDomain, err := ah.Usecase.GetSMSTemplateByName(r.Context(), templateName)
	if err != nil {
		handler.HandleErrors(w, r, err)
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, templateDomain)
}

// Get SMS template placeholders by name
// (GET /sms/templates/{templateName}/placeholders)
func (ah *APIHandler) GetSMSTemplatePlaceholdersByName(w http.ResponseWriter, r *http.Request, templateName string) {
	requiredClaims := map[string]any{"role": "user"}
	checkedAuthorization := ah.Authorizer.CheckIfAuthorised(w, r, requiredClaims)
	if !checkedAuthorization {
		return
	}
	templatePlaceholders, err := ah.Usecase.GetSMSPlaceholders(r.Context(), templateName)
	if err != nil {
		handler.HandleErrors(w, r, err)
		return
	}
	if len(templatePlaceholders) == 0 {
		handler.HandleErrors(w, r, domain.ErrorPlaceholderMissing{MissingPlaceholder: "No placeholders found in template"})
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, templatePlaceholders)
}
