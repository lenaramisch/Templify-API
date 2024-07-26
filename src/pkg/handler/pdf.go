package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

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
	var err error
	templateDomain, err := ah.usecase.GetPDFTemplateByName(templateName)
	if err != nil {
		handleError(res, req, http.StatusInternalServerError, "Error getting template")
		return
	}
	if templateDomain.TemplateStr == "" {
		handleError(res, req, http.StatusNotFound, fmt.Sprintf("Template with name %s not found", templateName))
		return
	}
	render.Status(req, http.StatusOK)
	render.JSON(res, req, templateDomain)
}

func (ah *APIHandler) GetFilledPDFTemplate(res http.ResponseWriter, req *http.Request) {
	templateName := chi.URLParam(req, "templateName")
	if templateName == "" {
		http.Error(res, "URL Param templateName empty", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, "Reading request body failed", http.StatusInternalServerError)
	}

	var values map[string]string
	if err := json.Unmarshal(body, &values); err != nil {
		http.Error(res, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	filledTemplate, err := ah.usecase.FillPDFTemplatePlaceholders(templateName, values)
	if err != nil {
		handleError(res, req, http.StatusInternalServerError, "Error filling template")
		return
	}
	render.Status(req, http.StatusOK)
	render.PlainText(res, req, filledTemplate)
}

func (ah *APIHandler) GeneratePDF(res http.ResponseWriter, req *http.Request) {
	templateName := chi.URLParam(req, "templateName")
	if templateName == "" {
		http.Error(res, "URL Param templateName empty", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, "Reading request body failed", http.StatusInternalServerError)
	}

	var values map[string]string
	if err := json.Unmarshal(body, &values); err != nil {
		http.Error(res, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	pdfBytes, err := ah.usecase.GeneratePDF(templateName, values)
	if err != nil {
		handleError(res, req, http.StatusInternalServerError, "Error generating PDF")
		return
	}

	render.Status(req, http.StatusOK)
	_, err = res.Write(pdfBytes)
	if err != nil {
		handleError(res, req, http.StatusInternalServerError, "Error responding with PDF")
		return
	}
}
