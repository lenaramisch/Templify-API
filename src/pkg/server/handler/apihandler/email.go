package apihandler

import (
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
// (POST /email)
func (ah *APIHandler) SendEmail(w http.ResponseWriter, r *http.Request) {
	var emailRequestAPI server.EmailSendRequest

	handler.ReadRequestBody(w, r, &emailRequestAPI)

	emailRequestDomain := &domain.EmailRequest{
		ToEmail:     emailRequestAPI.ToEmail,
		ToName:      emailRequestAPI.ToName,
		Subject:     emailRequestAPI.Subject,
		MessageBody: emailRequestAPI.Message,
	}
	err := ah.Usecase.SendRawEmail(emailRequestDomain)
	if err != nil {
		handler.HandleError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	render.Status(r, http.StatusOK)
	render.PlainText(w, r, "Email sent successfully")
}

// Send an Email with attachment
// (POST /email/attachments)
func (ah *APIHandler) SendEmailWithAttachment(w http.ResponseWriter, r *http.Request) {

	// Parse the multipart form, with a maximum memory of 32 MB for storing file parts in memory
	err := r.ParseMultipartForm(32 << 20) // 32MB
	if err != nil {
		handler.HandleError(w, r, http.StatusInternalServerError, "Error parsing multipart form")
		return
	}

	var emailRequest domain.EmailRequest
	emailRequest.ToEmail = r.FormValue("toEmail")
	emailRequest.ToName = r.FormValue("toName")
	emailRequest.Subject = r.FormValue("subject")
	message := r.FormValue("message")

	// TODO Do we really need to set this shouldBeSent value
	// or can we offer a Endpoint to see the rendered result?
	shouldBeSentString := r.FormValue("shouldBeSent")
	shouldBeSentBool := false
	if shouldBeSentString == "true" {
		shouldBeSentBool = true
	}

	emailRequest.ShouldBeSent = shouldBeSentBool
	attachmentInfo, err := handler.ReadMultipartFileAsBytes(r, w)
	if err != nil {
		http.Error(w, "Error with the attachment", http.StatusBadRequest)
		return
	}
	emailRequest.MessageBody = message
	emailRequest.AttachmentInfo = attachmentInfo

	if err := ah.Usecase.SendRawEmail(&emailRequest); err != nil {
		handler.HandleError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	render.Status(r, http.StatusOK)
	render.PlainText(w, r, "Email sent successfully")
}

// Get Template by Name
// (GET /templates/{templateName})
func (ah *APIHandler) GetTemplateByName(w http.ResponseWriter, r *http.Request, templateName string) {
	templateDomain := &domain.Template{}
	var err error
	templateDomain, err = ah.Usecase.GetEmailTemplateByName(templateName)
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
// (POST /templates/{templateName})
func (ah *APIHandler) AddNewTemplate(w http.ResponseWriter, r *http.Request, templateName string) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		handler.HandleError(w, r, http.StatusBadRequest, "Reading Request Body failed")
		return
	}
	MJMLString := string(body)

	err = ah.Usecase.AddEmailTemplate(templateName, MJMLString)
	if err != nil {
		handler.HandleError(w, r, http.StatusInternalServerError, fmt.Sprintf("Adding template with name %v failed", templateName))
		return
	}
	resultString := fmt.Sprintf("Added template with name %v", templateName)
	render.Status(r, http.StatusCreated)
	render.PlainText(w, r, resultString)
}

// Get Template Placeholders
// (GET /templates/{templateName}/placeholders)
func (ah *APIHandler) GetTemplatePlaceholdersByName(w http.ResponseWriter, r *http.Request, templateName string) {
	templatePlaceholders, err := ah.Usecase.GetEmailPlaceholders(templateName)
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
// (POST /templates/{templateName}/placeholders)
// TODO why is this expecting an INTEGER!?
func (ah *APIHandler) FillTemplate(w http.ResponseWriter, r *http.Request, templateName string) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Reading request body failed", http.StatusInternalServerError)
	}

	var templateFillRequest map[string]string

	if err := json.Unmarshal(body, &templateFillRequest); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	filledTemplate, err := ah.Usecase.GetFilledMJMLTemplate(templateName, templateFillRequest)
	if err != nil {
		handler.HandleError(w, r, http.StatusInternalServerError, "Error filling template")
		return
	}

	render.Status(r, http.StatusOK)
	render.PlainText(w, r, filledTemplate)
}

// Send a templated Email with attachment
// (POST /templates/{templateName}/placeholders/attachments)
func (ah *APIHandler) SendMJMLEmailWithAttachment(w http.ResponseWriter, r *http.Request, templateName string) {
	handler.FormToCapitalPlaceholders(r)
	attachmentInfo, err := handler.ReadMultipartFileAsBytes(r, w)
	if err != nil {
		http.Error(w, "Error with the attachment", http.StatusBadRequest)
		return
	}

	shouldBeSentString := r.FormValue("shouldBeSent")
	shouldBeSentBool := false
	if shouldBeSentString == "true" {
		shouldBeSentBool = true
	}

	// Fill emailRequest data
	domainEmailReq := &domain.EmailRequest{
		ToEmail:        r.FormValue("toEmail"),
		ToName:         r.FormValue("toName"),
		Subject:        r.FormValue("subject"),
		MessageBody:    "",
		ShouldBeSent:   shouldBeSentBool,
		AttachmentInfo: attachmentInfo,
	}

	if domainEmailReq.ToEmail == "" || domainEmailReq.ToName == "" || domainEmailReq.Subject == "" {
		http.Error(w, "Empty string content in either ToEmail, ToName, Subject", http.StatusBadRequest)
		return
	}
	// TODO read values from form
	var values map[string]string

	err = ah.Usecase.SendMJMLEmail(domainEmailReq, templateName, values)
	if err != nil {
		handler.HandleError(w, r, http.StatusInternalServerError, "Error sending mjml email")
		return
	}

	render.Status(r, http.StatusOK)
	render.PlainText(w, r, "sent!")
}
