package handler

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"unicode"

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

func handleError(res http.ResponseWriter, req *http.Request, statusCode int, message string) {
	render.Status(req, statusCode)
	render.PlainText(res, req, message)
}

func decodeJSONBody(res http.ResponseWriter, req *http.Request, dst any) error {
	err := json.NewDecoder(req.Body).Decode(dst)
	if err != nil {
		handleError(res, req, http.StatusBadRequest, "Invalid JSON: "+err.Error())
		return err
	}
	return nil
}

func (ah *APIHandler) SMSPostRequest(res http.ResponseWriter, req *http.Request) {
	var smsRequest SmsRequest

	if err := decodeJSONBody(res, req, &smsRequest); err != nil {
		return
	}

	if smsRequest.ToNumber == "" || smsRequest.MessageBody == "" {
		http.Error(res, "Empty string content in either ToNumber or MessageBody", http.StatusBadRequest)
		return
	}
	err := ah.usecase.SendSMS(smsRequest.ToNumber, smsRequest.MessageBody)
	if err != nil {
		handleError(res, req, http.StatusInternalServerError, err.Error())
		return
	}
	render.Status(req, http.StatusOK)
	render.PlainText(res, req, "SMS sent successfully")
}

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
	err := ah.usecase.SendEmail(&emailRequestDomain)
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
	templatePlaceholders, err := ah.usecase.GetTemplatePlaceholders(templateName)
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
	templateDomain, err = ah.usecase.GetTemplateByName(templateName)
	if err != nil {
		handleError(res, req, http.StatusInternalServerError, "Error getting template")
		return
	}
	if templateDomain.MJMLString == "" {
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

	var templateFillRequest domain.TemplateFillRequest

	if err := json.Unmarshal(body, &templateFillRequest); err != nil {
		http.Error(res, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	filledTemplate, err := ah.usecase.FillTemplatePlaceholders(templateName, &templateFillRequest)
	if err != nil {
		handleError(res, req, http.StatusInternalServerError, "Error filling template")
		return
	}

	render.Status(req, http.StatusOK)
	render.PlainText(res, req, filledTemplate)
}

func (ah *APIHandler) EmailPostRequestAttm(res http.ResponseWriter, req *http.Request) {
	var emailRequest domain.EmailRequestAttm

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

	shouldBeSentString := req.FormValue("shouldBeSent")
	shouldBeSentBool := false
	if shouldBeSentString == "true" {
		shouldBeSentBool = true
	}

	emailRequest.ShouldBeSent = shouldBeSentBool

	// Get the uploaded file
	file, handler, err := req.FormFile("file")
	if err != nil {
		handleError(res, req, http.StatusInternalServerError, "Error getting uploaded file")
		return
	}
	defer file.Close()

	// Read the file content into a byte slice
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		handleError(res, req, http.StatusInternalServerError, "Error reading file content")
		return
	}

	// Get file type string
	lastDotIndex := strings.LastIndex(handler.Filename, ".")
	if lastDotIndex == -1 {
		http.Error(res, "File name has to end with file extension, i.e. '.txt'", http.StatusBadRequest)
		return
	}

	fileTypeString := handler.Filename[lastDotIndex+1:]

	emailRequest.AttmContent = base64.StdEncoding.EncodeToString(fileBytes)
	emailRequest.FileName = handler.Filename
	emailRequest.FileType = fileTypeString

	if err := ah.usecase.SendEmailWithAttachment(message, &emailRequest); err != nil {
		handleError(res, req, http.StatusInternalServerError, err.Error())
		return
	}

	render.Status(req, http.StatusOK)
	render.PlainText(res, req, "Email sent successfully")
}

func (ah *APIHandler) PostTmplPlaceholdersAttm(res http.ResponseWriter, req *http.Request) {
	templateName := chi.URLParam(req, "templateName")
	if templateName == "" {
		http.Error(res, "URL Param templateName empty", http.StatusBadRequest)
		return
	}

	domainEmailRequest := &domain.EmailRequestAttm{}

	// Parse the multipart form, with a maximum memory of 32 MB for storing file parts in memory
	err := req.ParseMultipartForm(32 << 20) // 32MB
	if err != nil {
		handleError(res, req, http.StatusInternalServerError, "Error parsing multipart form")
		return
	}

	form := req.MultipartForm
	domainEmailRequest.Placeholders = make(map[string]string)
	for key, values := range form.Value {
		if len(key) > 0 && unicode.IsUpper(rune(key[0])) {
			domainEmailRequest.Placeholders[key] = values[0]
		}
	}

	file, handler, err := req.FormFile("file")
	if err != nil {
		handleError(res, req, http.StatusInternalServerError, "Error getting uploaded file")
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		handleError(res, req, http.StatusInternalServerError, "Error reading file content")
		return
	}

	// Get file type string
	lastDotIndex := strings.LastIndex(handler.Filename, ".")
	if lastDotIndex == -1 {
		http.Error(res, "File name has to end with file extension, i.e. '.txt'", http.StatusBadRequest)
		return
	}

	fileTypeString := handler.Filename[lastDotIndex+1:]

	shouldBeSentString := req.FormValue("shouldBeSent")
	shouldBeSentBool := false
	if shouldBeSentString == "true" {
		shouldBeSentBool = true
	}

	// Fill emailRequest data
	domainEmailRequest.FileType = fileTypeString
	domainEmailRequest.AttmContent = base64.StdEncoding.EncodeToString(fileBytes)
	domainEmailRequest.FileName = handler.Filename
	domainEmailRequest.ToEmail = req.FormValue("toEmail")
	domainEmailRequest.ToName = req.FormValue("toName")
	domainEmailRequest.Subject = req.FormValue("subject")
	domainEmailRequest.ShouldBeSent = shouldBeSentBool

	if domainEmailRequest.ToEmail == "" || domainEmailRequest.ToName == "" || domainEmailRequest.Subject == "" {
		http.Error(res, "Empty string content in either ToEmail, ToName, Subject", http.StatusBadRequest)
		return
	}

	filledTemplate, err := ah.usecase.FillTemplatePlaceholdersAttm(templateName, domainEmailRequest)
	if err != nil {
		handleError(res, req, http.StatusInternalServerError, "Error filling template")
		return
	}

	render.Status(req, http.StatusOK)
	render.PlainText(res, req, filledTemplate)
}

func (ah *APIHandler) PDFTemplPostReq(res http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		handleError(res, req, http.StatusBadRequest, "Reading Request Body failed")
		return
	}
	typstString := string(body)

	fmt.Println("Typst String: ", typstString)
	templateName := chi.URLParam(req, "templateName")
	if templateName == "" {
		http.Error(res, "URL Param templateName empty", http.StatusBadRequest)
		return
	}

	err = ah.usecase.AddPDFTemplate(templateName, typstString)
	if err != nil {
		handleError(res, req, http.StatusInternalServerError, fmt.Sprintf("Adding template with name %v failed", templateName))
		return
	}
	resultString := fmt.Sprintf("Added template with name %v", templateName)
	render.Status(req, http.StatusCreated)
	render.PlainText(res, req, resultString)
}

func (ah *APIHandler) GetPDFTemplByName(res http.ResponseWriter, req *http.Request) {
	templateName := chi.URLParam(req, "templateName")
	if templateName == "" {
		http.Error(res, "URL Param templateName empty", http.StatusBadRequest)
		return
	}
	templateDomain := &domain.PDFTemplate{}
	var err error
	templateDomain, err = ah.usecase.GetPDFTemplateByName(templateName)
	if err != nil {
		handleError(res, req, http.StatusInternalServerError, "Error getting template")
		return
	}
	if templateDomain.TypstString == "" {
		handleError(res, req, http.StatusNotFound, fmt.Sprintf("Template with name %s not found", templateName))
		return
	}
	render.Status(req, http.StatusOK)
	render.JSON(res, req, templateDomain)
}
