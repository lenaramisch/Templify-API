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

// Create a new workflow
// (POST /workflow/{workflowName})
func (ah *APIHandler) CreateWorkflow(w http.ResponseWriter, r *http.Request, workflowName string) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		handler.HandleError(w, r, http.StatusBadRequest, "Reading Request Body failed")
		return
	}

	var addWorkflowRequest server.WorkflowCreateRequest
	if err := json.Unmarshal(body, &addWorkflowRequest); err != nil {
		handler.HandleError(w, r, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	// Map the DTO to the domain model
	workflowDomain := &domain.WorkflowCreateRequest{
		Name:              workflowName,
		EmailSubject:      addWorkflowRequest.EmailSubject,
		EmailTemplateName: addWorkflowRequest.EmailTemplateName,
		StaticAttachments: addWorkflowRequest.StaticAttachmentNames,
		TemplatedPDFs:     addWorkflowRequest.TemplatedAttachmentNames,
	}

	err = ah.Usecase.AddWorkflow(r.Context(), workflowDomain)
	if err != nil {
		handler.HandleError(w, r, http.StatusInternalServerError, fmt.Sprintf("Adding workflow with name %v failed", workflowName))
		return
	}
	resultString := fmt.Sprintf("Added workflow with name %v", workflowName)
	render.Status(r, http.StatusCreated)
	render.PlainText(w, r, resultString)
}

// Get a workflow by name
// (GET /workflow/{workflowName})
func (ah *APIHandler) GetWorkflowByName(w http.ResponseWriter, r *http.Request, workflowName string) {
	workflowDomain, err := ah.Usecase.GetWorkflowByName(r.Context(), workflowName)
	if err != nil {
		handler.HandleError(w, r, http.StatusInternalServerError, "Error getting workflow")
		return
	}
	if workflowDomain.Name == "" {
		handler.HandleError(w, r, http.StatusNotFound, fmt.Sprintf("Workflow with name %s not found", workflowName))
		return
	}

	// Map the domain model to the DTO
	emailTemplateInfo := &server.TemplateInfo{
		Placeholders: workflowDomain.EmailTemplate.Placeholders,
		TemplateName: workflowDomain.EmailTemplate.TemplateName,
	}

	// range over the pdf templates and map them to the DTO
	var pdfTemplates []server.TemplateInfo
	if workflowDomain.PDFTemplates != nil {
		for _, pdfTemplate := range workflowDomain.PDFTemplates {
			pdfTemplateInfo := server.TemplateInfo{
				Placeholders: pdfTemplate.Placeholders,
				TemplateName: pdfTemplate.TemplateName,
			}
			pdfTemplates = append(pdfTemplates, pdfTemplateInfo)
		}
	}

	getWorkflowResponse := &server.GetWorkflowResponse{
		Name:              workflowDomain.Name,
		EmailSubject:      workflowDomain.EmailSubject,
		EmailTemplate:     *emailTemplateInfo,
		PdfTemplates:      pdfTemplates,
		StaticAttachments: workflowDomain.StaticAttachments,
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, getWorkflowResponse)
}

// Use a workflow by name
// (POST /workflow/{workflowName}/send)
func (ah *APIHandler) UseWorkflow(w http.ResponseWriter, r *http.Request, workflowName string) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		handler.HandleError(w, r, http.StatusBadRequest, "Reading Request Body failed")
		return
	}

	var useWorkflowRequest server.WorkflowSendRequest
	if err := json.Unmarshal(body, &useWorkflowRequest); err != nil {
		handler.HandleError(w, r, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	// Map the DTO to the domain model
	useWorkflowRequestDomain := &domain.WorkflowUseRequest{
		Name: workflowName,
		EmailTemplate: domain.TemplateToFill{
			Placeholders: *useWorkflowRequest.EmailTemplate.Placeholders,
			TemplateName: *useWorkflowRequest.EmailTemplate.TemplateName},
		ToEmail: useWorkflowRequest.ToEmail,
		ToName:  useWorkflowRequest.ToName,
	}

	//TODO How to handle multiple PDF templates?
	if useWorkflowRequest.PdfTemplates != nil {
		// range over the pdf templates and map them to the domain model
		for _, pdfTemplate := range *useWorkflowRequest.PdfTemplates {
			useWorkflowRequestDomain.PdfTemplates = append(useWorkflowRequestDomain.PdfTemplates, domain.TemplateToFill{
				Placeholders: *pdfTemplate.Placeholders,
				TemplateName: *pdfTemplate.TemplateName,
			})
		}
	}

	err = ah.Usecase.UseWorkflow(r.Context(), useWorkflowRequestDomain)
	if err != nil {
		ah.log.With("error", err).Debug("Error using workflow")
		handler.HandleError(w, r, http.StatusInternalServerError, fmt.Sprintf("Using workflow with name %v failed", workflowName))
		return
	}
	resultString := fmt.Sprintf("Used workflow with name %v", workflowName)
	render.Status(r, http.StatusCreated)
	render.PlainText(w, r, resultString)
}
