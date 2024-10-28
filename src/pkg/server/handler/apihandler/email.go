package apihandler

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	domain "templify/pkg/domain/model"
	server "templify/pkg/server/generated"
	"templify/pkg/server/handler"

	"github.com/go-chi/render"
)

// Send an Email with custom text
// (POST /email/basic/send)
func (ah *APIHandler) SendEmail(w http.ResponseWriter, r *http.Request) {
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
				handler.HandleError(w, r, http.StatusBadRequest, "Invalid base64 string")
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
		handler.HandleError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	render.Status(r, http.StatusOK)
	render.PlainText(w, r, "Email sent successfully")
}

// Send an Email with template
// (POST /email/templates/{templateName}/send)
func (ah *APIHandler) SendTemplatedEmail(w http.ResponseWriter, r *http.Request, templateName string) {
	requiredClaims := map[string]any{"role": "user"}
	checkedAuthorization := handler.CheckIfAuthorised(w, r, requiredClaims)
	if !checkedAuthorization {
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Reading request body failed", http.StatusInternalServerError)
		return
	}

	// Create EmailTemplateSendRequest
	var emailTemplateSendReq server.EmailTemplateFillSendRequest

	// Read the request body and unmarshal it into emailTemplateSendReq
	if err := json.Unmarshal(body, &emailTemplateSendReq); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	var attachmentInfo []domain.AttachmentInfo
	// Check if Attachments are present and create AttachmentInfo
	if emailTemplateSendReq.Attachments != nil {
		for _, attachment := range *emailTemplateSendReq.Attachments {
			bytes, err := base64.StdEncoding.DecodeString(attachment.AttachmentContent)
			if err != nil {
				handler.HandleError(w, r, http.StatusBadRequest, "Invalid base64 string")
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
		handler.HandleError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	render.Status(r, http.StatusOK)
	render.PlainText(w, r, "Email sent successfully")
}

// Get Template by Name
// (GET /email/templates/{templateName})
func (ah *APIHandler) GetTemplateByName(w http.ResponseWriter, r *http.Request, templateName string) {
	templateDomain, err := ah.Usecase.GetEmailTemplateByName(r.Context(), templateName)
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

// Add new template
// (POST /email/templates/{templateName})
func (ah *APIHandler) AddNewTemplate(w http.ResponseWriter, r *http.Request, templateName string) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		handler.HandleError(w, r, http.StatusBadRequest, "Reading Request Body failed")
		return
	}

	var addTemplateRequest server.EmailTemplatePostRequest
	if err := json.Unmarshal(body, &addTemplateRequest); err != nil {
		handler.HandleError(w, r, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	var templateDomain = &domain.Template{
		Name:        templateName,
		TemplateStr: addTemplateRequest.TemplateString,
		IsMJML:      addTemplateRequest.IsMJML,
	}

	err = ah.Usecase.AddEmailTemplate(r.Context(), templateDomain)
	if err != nil {
		handler.HandleError(w, r, http.StatusInternalServerError, fmt.Sprintf("Adding template with name %v failed", templateName))
		return
	}
	resultString := fmt.Sprintf("Added template with name %v", templateName)
	render.Status(r, http.StatusCreated)
	render.PlainText(w, r, resultString)
}

// Get Template Placeholders
// (GET /email/templates/{templateName}/placeholders)
func (ah *APIHandler) GetTemplatePlaceholdersByName(w http.ResponseWriter, r *http.Request, templateName string) {
	templatePlaceholders, err := ah.Usecase.GetEmailPlaceholders(r.Context(), templateName)
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

// Fill placeholders of template
// (POST /email/templates/{templateName}/fill)
func (ah *APIHandler) FillTemplate(w http.ResponseWriter, r *http.Request, templateName string) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Reading request body failed", http.StatusInternalServerError)
	}

	var templateFillRequest server.EmailTemplateFillSendRequest

	if err := json.Unmarshal(body, &templateFillRequest); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	filledTemplate, err := ah.Usecase.GetFilledTemplateString(r.Context(), templateName, templateFillRequest.Placeholders)
	if err != nil {
		handler.HandleError(w, r, http.StatusInternalServerError, "Error filling template")
		return
	}

	render.Status(r, http.StatusOK)
	render.PlainText(w, r, filledTemplate)
}
