package apihandler

import (
	"fmt"
	"net/http"
	domain "templify/pkg/domain/model"
	server "templify/pkg/server/generated"
	"templify/pkg/server/handler"

	"github.com/go-chi/render"
)

// Create a new workflow
// (POST /workflow/{workflowName})
func (ah *APIHandler) CreateWorkflow(w http.ResponseWriter, r *http.Request, workflowName string) {
	requiredClaims := map[string]any{"role": "user"}
	checkedAuthorization := ah.Authorizer.CheckIfAuthorised(w, r, requiredClaims)
	if !checkedAuthorization {
		return
	}

	var addWorkflowRequest server.WorkflowCreateRequest
	err := handler.ReadRequestBody(w, r, &addWorkflowRequest)
	if err != nil {
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
		handler.HandleErrors(w, r, err)
		return
	}

	render.Status(r, http.StatusCreated)
	render.PlainText(w, r, fmt.Sprintf("Added workflow with name %v", workflowName))
}

// Get a workflow by name
// (GET /workflow/{workflowName})
func (ah *APIHandler) GetWorkflowByName(w http.ResponseWriter, r *http.Request, workflowName string) {
	requiredClaims := map[string]any{"role": "user"}
	checkedAuthorization := ah.Authorizer.CheckIfAuthorised(w, r, requiredClaims)
	if !checkedAuthorization {
		return
	}
	workflowDomain, err := ah.Usecase.GetWorkflowByName(r.Context(), workflowName)
	if err != nil {
		handler.HandleErrors(w, r, err)
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
	requiredClaims := map[string]any{"role": "user"}
	checkedAuthorization := ah.Authorizer.CheckIfAuthorised(w, r, requiredClaims)
	if !checkedAuthorization {
		return
	}

	var useWorkflowRequest server.WorkflowSendRequest
	err := handler.ReadRequestBody(w, r, &useWorkflowRequest)
	if err != nil {
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
		handler.HandleErrors(w, r, err)
		return
	}

	render.Status(r, http.StatusCreated)
	render.PlainText(w, r, fmt.Sprintf("Used workflow with name %v", workflowName))
}
