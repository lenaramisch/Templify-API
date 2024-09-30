package apihandler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	server "templify/pkg/server/generated"
	"templify/pkg/server/handler"

	"github.com/go-chi/render"
)

func (ah *APIHandler) GetPDFTemplatePlaceholdersByName(w http.ResponseWriter, r *http.Request, templateName string) {

	if templateName == "" {
		http.Error(w, "URL Param templateName empty", http.StatusBadRequest)
		return
	}

	templatePlaceholders, err := ah.Usecase.GetPDFPlaceholders(templateName)
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

func (ah *APIHandler) AddNewPDFTemplate(w http.ResponseWriter, r *http.Request, templateName string) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		handler.HandleError(w, r, http.StatusBadRequest, "Reading Request Body failed")
		return
	}

	var addTemplateRequest server.PDFTemplatePostRequest
	if err := json.Unmarshal(body, &addTemplateRequest); err != nil {
		handler.HandleError(w, r, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	if templateName == "" {
		http.Error(w, "URL Param templateName empty", http.StatusBadRequest)
		return
	}

	err = ah.Usecase.AddPDFTemplate(templateName, addTemplateRequest.TemplateString)
	if err != nil {
		handler.HandleError(w, r, http.StatusInternalServerError, fmt.Sprintf("Adding template with name %v failed", templateName))
		return
	}
	resultString := fmt.Sprintf("Added template with name %v", templateName)
	render.Status(r, http.StatusCreated)
	render.PlainText(w, r, resultString)
}

func (ah *APIHandler) GetPDFTemplateByName(w http.ResponseWriter, r *http.Request, templateName string) {
	if templateName == "" {
		http.Error(w, "URL Param templateName empty", http.StatusBadRequest)
		return
	}
	var err error
	templateDomain, err := ah.Usecase.GetPDFTemplateByName(templateName)
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

func (ah *APIHandler) FillPDFTemplate(w http.ResponseWriter, r *http.Request, templateName string) {
	if templateName == "" {
		http.Error(w, "URL Param templateName empty", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Reading request body failed", http.StatusInternalServerError)
	}

	var pdfFillReq server.TemplateFillRequest
	if err := json.Unmarshal(body, &pdfFillReq); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	pdfBytes, err := ah.Usecase.GeneratePDF(templateName, pdfFillReq.Placeholders)
	if err != nil {
		handler.HandleError(w, r, http.StatusInternalServerError, "Error generating PDF")
		return
	}

	render.Status(r, http.StatusOK)
	_, err = w.Write(pdfBytes)
	if err != nil {
		handler.HandleError(w, r, http.StatusInternalServerError, "Error responding with PDF")
		return
	}
}
