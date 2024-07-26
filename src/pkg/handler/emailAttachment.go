package handler

import (
	"net/http"

	"example.SMSService.com/pkg/domain"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func (ah *APIHandler) SendEmailWithAttachment(res http.ResponseWriter, req *http.Request) {
	var emailRequest domain.EmailRequest

	// Parse the multipart form, with a maximum memory of 32 MB for storing file parts in memory
	err := req.ParseMultipartForm(32 << 20) // 32MB
	if err != nil {
		handleError(res, req, http.StatusInternalServerError, "Error parsing multipart form")
		return
	}

	emailRequest.ToEmail = req.FormValue("toEmail")
	emailRequest.ToName = req.FormValue("toName")
	emailRequest.Subject = req.FormValue("subject")
	message := req.FormValue("message")

	if emailRequest.ToEmail == "" || emailRequest.ToName == "" || emailRequest.Subject == "" || message == "" {
		http.Error(res, "Empty string content in either ToEmail, ToName, Subject or MessageBody", http.StatusBadRequest)
		return
	}

	// TODO Do we really need to set this shouldBeSent value
	// or can we offer a Endpoint to see the rendered result?
	shouldBeSentString := req.FormValue("shouldBeSent")
	shouldBeSentBool := false
	if shouldBeSentString == "true" {
		shouldBeSentBool = true
	}

	emailRequest.ShouldBeSent = shouldBeSentBool
	attachmentInfo, err := readMultipartFileAsBytes(req, res)
	if err != nil {
		http.Error(res, "Error with the attachment", http.StatusBadRequest)
		return
	}
	emailRequest.MessageBody = message
	emailRequest.AttachmentInfo = attachmentInfo

	if err := ah.usecase.SendRawEmail(&emailRequest); err != nil {
		handleError(res, req, http.StatusInternalServerError, err.Error())
		return
	}

	render.Status(req, http.StatusOK)
	render.PlainText(res, req, "Email sent successfully")
}

func (ah *APIHandler) SendMJMLWithAttachment(res http.ResponseWriter, req *http.Request) {
	templateName := chi.URLParam(req, "templateName")
	if templateName == "" {
		http.Error(res, "URL Param templateName empty", http.StatusBadRequest)
		return
	}

	formToCapitalPlaceholders(req)
	attachmentInfo, err := readMultipartFileAsBytes(req, res)
	if err != nil {
		http.Error(res, "Error with the attachment", http.StatusBadRequest)
		return
	}

	shouldBeSentString := req.FormValue("shouldBeSent")
	shouldBeSentBool := false
	if shouldBeSentString == "true" {
		shouldBeSentBool = true
	}

	// Fill emailRequest data
	domainEmailReq := &domain.EmailRequest{
		ToEmail:        req.FormValue("toEmail"),
		ToName:         req.FormValue("toName"),
		Subject:        req.FormValue("subject"),
		MessageBody:    "",
		ShouldBeSent:   shouldBeSentBool,
		AttachmentInfo: attachmentInfo,
	}

	if domainEmailReq.ToEmail == "" || domainEmailReq.ToName == "" || domainEmailReq.Subject == "" {
		http.Error(res, "Empty string content in either ToEmail, ToName, Subject", http.StatusBadRequest)
		return
	}
	// TODO read values from form
	var values map[string]string

	err = ah.usecase.SendMJMLEmail(domainEmailReq, templateName, values)
	if err != nil {
		handleError(res, req, http.StatusInternalServerError, "Error sending mjml email")
		return
	}

	render.Status(req, http.StatusOK)
	render.PlainText(res, req, "sent!")
}
