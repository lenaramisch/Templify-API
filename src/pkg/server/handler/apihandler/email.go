package apihandler

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	domain "templify/pkg/domain/model"
	server "templify/pkg/server/generated"
	"templify/pkg/server/handler"

	"github.com/go-chi/render"
)

// Send an Email with custom text
// (POST /email/basic/send)
func (ah *APIHandler) SendEmail(w http.ResponseWriter, r *http.Request) {
	requiredClaims := map[string]any{"role": "user"}
	checkedAuthorization := ah.Authorizer.CheckIfAuthorised(w, r, requiredClaims)
	if !checkedAuthorization {
		return
	}
	var emailRequestAPI server.EmailSendRequest
	err := handler.ReadRequestBody(w, r, &emailRequestAPI)
	if err != nil {
		ah.log.Debug("Error reading request body")
		return
	}

	var attachmentInfo []domain.AttachmentInfo
	// Check if Attachments are present and create AttachmentInfo
	if emailRequestAPI.Attachments != nil {
		for _, attachment := range *emailRequestAPI.Attachments {
			// convert base64 string attachmentContent into []byte
			bytes, err := base64.StdEncoding.DecodeString(attachment.AttachmentContent)
			if err != nil {
				handler.HandleErrors(w, r, errors.New("invalid base64 string"))
			}

			attachmentInfo = append(attachmentInfo, domain.AttachmentInfo{
				FileExtension: attachment.AttachmentExtension,
				FileName:      attachment.AttachmentName,
				FileBytes:     bytes,
			})
		}
	}

	emailRequestDomain := &domain.EmailRequest{
		AttachmentInfo: attachmentInfo,
		ToEmail:        emailRequestAPI.ToEmail,
		ToName:         emailRequestAPI.ToName,
		Subject:        emailRequestAPI.Subject,
		MessageBody:    emailRequestAPI.Message,
	}

	err = ah.Usecase.SendRawEmail(emailRequestDomain)
	if err != nil {
		handler.HandleErrors(w, r, err)
		return
	}
	render.Status(r, http.StatusOK)
	render.PlainText(w, r, "Email sent successfully")
}

// Send an Email with template
// (POST /email/templates/{templateName}/send)
func (ah *APIHandler) SendTemplatedEmail(w http.ResponseWriter, r *http.Request, templateName string) {
	requiredClaims := map[string]any{"role": "user"}
	checkedAuthorization := ah.Authorizer.CheckIfAuthorised(w, r, requiredClaims)
	if !checkedAuthorization {
		return
	}
	var emailTemplateSendReq server.EmailTemplateFillSendRequest
	err := handler.ReadRequestBody(w, r, &emailTemplateSendReq)
	if err != nil {
		return
	}

	var attachmentInfo []domain.AttachmentInfo
	// Check if Attachments are present and create AttachmentInfo
	if emailTemplateSendReq.Attachments != nil {
		for _, attachment := range *emailTemplateSendReq.Attachments {
			bytes, err := base64.StdEncoding.DecodeString(attachment.AttachmentContent)
			if err != nil {
				handler.HandleErrors(w, r, errors.New("invalid base64 string"))
			}
			attachmentInfo = append(attachmentInfo, domain.AttachmentInfo{
				FileExtension: attachment.AttachmentExtension,
				FileName:      attachment.AttachmentName,
				FileBytes:     bytes,
			})
		}
	}

	// Create emailRequestDomain
	emailRequestDomain := &domain.EmailTemplateSendRequest{
		AttachmentInfo: attachmentInfo,
		ToEmail:        emailTemplateSendReq.ToEmail,
		ToName:         emailTemplateSendReq.ToName,
		Subject:        emailTemplateSendReq.Subject,
		Placeholders:   emailTemplateSendReq.Placeholders,
		TemplateName:   templateName,
	}

	// Call Usecase to send email
	err = ah.Usecase.SendTemplatedEmail(r.Context(), emailRequestDomain)
	if err != nil {
		handler.HandleErrors(w, r, err)
		return
	}
	render.Status(r, http.StatusOK)
	render.PlainText(w, r, "Email sent successfully")
}

// Get Template by Name
// (GET /email/templates/{templateName})
func (ah *APIHandler) GetTemplateByName(w http.ResponseWriter, r *http.Request, templateName string) {
	requiredClaims := map[string]any{"role": "user"}
	checkedAuthorization := ah.Authorizer.CheckIfAuthorised(w, r, requiredClaims)
	if !checkedAuthorization {
		return
	}
	templateDomain, err := ah.Usecase.GetEmailTemplateByName(r.Context(), templateName)
	if err != nil {
		handler.HandleErrors(w, r, err)
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, templateDomain)
}

// Add new template
// (POST /email/templates/{templateName})
func (ah *APIHandler) AddNewTemplate(w http.ResponseWriter, r *http.Request, templateName string) {
	requiredClaims := map[string]any{"role": "user"}
	checkedAuthorization := ah.Authorizer.CheckIfAuthorised(w, r, requiredClaims)
	if !checkedAuthorization {
		return
	}
	var addTemplateRequest server.EmailTemplatePostRequest
	err := handler.ReadRequestBody(w, r, &addTemplateRequest)
	if err != nil {
		return
	}

	var templateDomain = &domain.Template{
		Name:        templateName,
		TemplateStr: addTemplateRequest.TemplateString,
		IsMJML:      addTemplateRequest.IsMJML,
	}

	err = ah.Usecase.AddEmailTemplate(r.Context(), templateDomain)
	if err != nil {
		handler.HandleErrors(w, r, err)
		return
	}
	render.Status(r, http.StatusCreated)
	render.PlainText(w, r, fmt.Sprintf("Added template with name %v", templateName))
}

// Get Template Placeholders
// (GET /email/templates/{templateName}/placeholders)
func (ah *APIHandler) GetTemplatePlaceholdersByName(w http.ResponseWriter, r *http.Request, templateName string) {
	requiredClaims := map[string]any{"role": "user"}
	checkedAuthorization := ah.Authorizer.CheckIfAuthorised(w, r, requiredClaims)
	if !checkedAuthorization {
		return
	}
	templatePlaceholders, err := ah.Usecase.GetEmailPlaceholders(r.Context(), templateName)
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

// Fill placeholders of template
// (POST /email/templates/{templateName}/fill)
func (ah *APIHandler) FillTemplate(w http.ResponseWriter, r *http.Request, templateName string) {
	requiredClaims := map[string]any{"role": "user"}
	checkedAuthorization := ah.Authorizer.CheckIfAuthorised(w, r, requiredClaims)
	if !checkedAuthorization {
		return
	}

	var templateFillRequest server.EmailTemplateFillSendRequest
	err := handler.ReadRequestBody(w, r, &templateFillRequest)
	if err != nil {
		return
	}

	filledTemplate, err := ah.Usecase.GetFilledTemplateString(r.Context(), templateName, templateFillRequest.Placeholders)
	if err != nil {
		handler.HandleErrors(w, r, err)
	}

	render.Status(r, http.StatusOK)
	render.PlainText(w, r, filledTemplate)
}
