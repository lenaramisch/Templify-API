package apihandler

import (
	"fmt"
	"net/http"
	domain "templify/pkg/domain/model"
	server "templify/pkg/server/generated"
	"templify/pkg/server/handler"

	"github.com/go-chi/render"
)

// Get all placehholders of a PDF template
// (GET /pdf/templates/{templateName}/placeholders)
func (ah *APIHandler) GetPDFTemplatePlaceholdersByName(w http.ResponseWriter, r *http.Request, templateName string) {
	requiredClaims := map[string]any{"role": "user"}
	checkedAuthorization := ah.Authorizer.CheckIfAuthorised(w, r, requiredClaims)
	if !checkedAuthorization {
		return
	}

	templatePlaceholders, err := ah.Usecase.GetPDFPlaceholders(r.Context(), templateName)
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

// Add new PDF template
// (POST /pdf/templates/{templateName})
func (ah *APIHandler) AddNewPDFTemplate(w http.ResponseWriter, r *http.Request, templateName string) {
	requiredClaims := map[string]any{"role": "user"}
	checkedAuthorization := ah.Authorizer.CheckIfAuthorised(w, r, requiredClaims)
	if !checkedAuthorization {
		return
	}

	var addTemplateRequest server.PDFTemplatePostRequest
	err := handler.ReadRequestBody(w, r, &addTemplateRequest)
	if err != nil {
		return
	}

	templateDomain := &domain.Template{
		Name:        templateName,
		TemplateStr: addTemplateRequest.TemplateString,
	}

	err = ah.Usecase.AddPDFTemplate(r.Context(), templateDomain)
	if err != nil {
		handler.HandleErrors(w, r, err)
		return
	}
	render.Status(r, http.StatusCreated)
	render.PlainText(w, r, fmt.Sprintf("Added template with name %v", templateName))
}

// Get PDF template by name
// (GET /pdf/templates/{templateName})
func (ah *APIHandler) GetPDFTemplateByName(w http.ResponseWriter, r *http.Request, templateName string) {
	requiredClaims := map[string]any{"role": "user"}
	checkedAuthorization := ah.Authorizer.CheckIfAuthorised(w, r, requiredClaims)
	if !checkedAuthorization {
		return
	}
	templateDomain, err := ah.Usecase.GetPDFTemplateByName(r.Context(), templateName)
	if err != nil {
		handler.HandleErrors(w, r, err)
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, templateDomain)
}

// Fill placeholders of PDF template
// (POST /pdf/templates/{templateName}/fill)
func (ah *APIHandler) FillPDFTemplate(w http.ResponseWriter, r *http.Request, templateName string) {
	requiredClaims := map[string]any{"role": "user"}
	checkedAuthorization := ah.Authorizer.CheckIfAuthorised(w, r, requiredClaims)
	if !checkedAuthorization {
		return
	}

	var pdfFillReq server.TemplateFillRequest
	err := handler.ReadRequestBody(w, r, &pdfFillReq)
	if err != nil {
		return
	}

	pdfBytes, err := ah.Usecase.GeneratePDF(r.Context(), templateName, pdfFillReq.Placeholders)
	if err != nil {
		handler.HandleErrors(w, r, err)
		return
	}

	_, err = w.Write(pdfBytes)
	if err != nil {
		handler.HandleErrors(w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
}
