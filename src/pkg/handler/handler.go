package handler

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

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
	render.Status(req, http.StatusCreated)
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
		http.Error(res, fmt.Sprintf("Getting placeholders for template %s failed", templateName), http.StatusInternalServerError)
		return
	}
	if len(templatePlaceholders) == 0 {
		http.Error(res, fmt.Sprintf("No placeholders for template %s found", templateName), http.StatusNotFound)
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
	if templateDomain.MJMLString == "" {
		http.Error(res, fmt.Sprintf("Template with name %s not found", templateName), http.StatusNotFound)
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

	templateFillRequest := TemplateFillRequest{}

	if err := json.Unmarshal(body, &templateFillRequest); err != nil {
		http.Error(res, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	filledTemplate, err := ah.usecase.FillTemplatePlaceholders(templateName, templateFillRequest.ShouldBeSent, templateFillRequest.ToEmail, templateFillRequest.ToName, templateFillRequest.Subject, templateFillRequest.Placeholders)
	if err != nil {
		http.Error(res, "Error filling template", http.StatusInternalServerError)
		return
	}

	render.Status(req, http.StatusOK)
	render.PlainText(res, req, filledTemplate)
}

func (ah *APIHandler) EmailPostRequestAttachment(res http.ResponseWriter, req *http.Request) {
	var emailRequest EmailAttachmentRequest

	// Parse the multipart form, with a maximum memory of 32 MB for storing file parts in memory
	err := req.ParseMultipartForm(32 << 20) // 32MB
	if err != nil {
		fmt.Print("Error trying to parse multipart form")
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	emailRequest.ToEmail = req.FormValue("toEmail")
	emailRequest.ToName = req.FormValue("toName")
	emailRequest.Subject = req.FormValue("subject")
	emailRequest.MessageBody = req.FormValue("message")

	if emailRequest.ToEmail == "" || emailRequest.ToName == "" || emailRequest.Subject == "" || emailRequest.MessageBody == "" {
		http.Error(res, "Empty string content in either ToEmail, ToName, Subject or MessageBody", http.StatusBadRequest)
		return
	}

	// Get the uploaded file
	file, handler, err := req.FormFile("file")
	if err != nil {
		fmt.Print("Error getting uploaded file")
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Read the file content into a byte slice
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		fmt.Print("Error reading file content into byte slice")
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	// Encode the byte slice to a base64 string
	base64Str := base64.StdEncoding.EncodeToString(fileBytes)

	emailRequest.AttachmentContent = base64Str
	emailRequest.FileName = handler.Filename

	lastDotIndex := strings.LastIndex(handler.Filename, ".")
	if lastDotIndex == -1 {
		http.Error(res, "File name has to end with file extension, i.e. '.txt'", http.StatusBadRequest)
		return
	}

	fileTypeString := handler.Filename[lastDotIndex+1:]

	emailRequest.FileType = fileTypeString

	// Print the base64 encoded string
	fmt.Fprintf(res, "Base64 Encoded File: %s\n", base64Str)
	err = ah.usecase.SendEmailWithAttachment(emailRequest.ToEmail, emailRequest.ToName, emailRequest.Subject, emailRequest.MessageBody, emailRequest.AttachmentContent, emailRequest.FileName, emailRequest.FileType)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	render.Status(req, http.StatusOK)
	render.PlainText(res, req, "Email sent successfully")
}
