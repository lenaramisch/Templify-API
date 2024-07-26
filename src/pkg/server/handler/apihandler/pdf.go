package apihandler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"templify/pkg/server/handler"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func (ah *APIHandler) PDFTemplPostReq(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		handler.HandleError(w, r, http.StatusBadRequest, "Reading Request Body failed")
		return
	}
	typstString := string(body)

	fmt.Println("Typst String: ", typstString)
	templateName := chi.URLParam(r, "templateName")
	if templateName == "" {
		http.Error(w, "URL Param templateName empty", http.StatusBadRequest)
		return
	}

	err = ah.Usecase.AddPDFTemplate(templateName, typstString)
	if err != nil {
		handler.HandleError(w, r, http.StatusInternalServerError, fmt.Sprintf("Adding template with name %v failed", templateName))
		return
	}
	resultString := fmt.Sprintf("Added template with name %v", templateName)
	render.Status(r, http.StatusCreated)
	render.PlainText(w, r, resultString)
}

func (ah *APIHandler) GetPDFTemplByName(w http.ResponseWriter, r *http.Request) {
	templateName := chi.URLParam(r, "templateName")
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

func (ah *APIHandler) GetFilledPDFTemplate(w http.ResponseWriter, r *http.Request) {
	templateName := chi.URLParam(r, "templateName")
	if templateName == "" {
		http.Error(w, "URL Param templateName empty", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Reading request body failed", http.StatusInternalServerError)
	}

	var values map[string]string
	if err := json.Unmarshal(body, &values); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	filledTemplate, err := ah.Usecase.FillPDFTemplatePlaceholders(templateName, values)
	if err != nil {
		handler.HandleError(w, r, http.StatusInternalServerError, "Error filling template")
		return
	}
	render.Status(r, http.StatusOK)
	render.PlainText(w, r, filledTemplate)
}

func (ah *APIHandler) GeneratePDF(w http.ResponseWriter, r *http.Request) {
	templateName := chi.URLParam(r, "templateName")
	if templateName == "" {
		http.Error(w, "URL Param templateName empty", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Reading request body failed", http.StatusInternalServerError)
	}

	var values map[string]string
	if err := json.Unmarshal(body, &values); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	pdfBytes, err := ah.Usecase.GeneratePDF(templateName, values)
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
