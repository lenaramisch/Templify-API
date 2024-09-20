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
		Name:                workflowName,
		EmailSubject:        addWorkflowRequest.EmailSubject,
		EmailTemplateName:   addWorkflowRequest.EmailTemplate.EmailTemplateName,
		EmailTemplateString: addWorkflowRequest.EmailTemplate.EmailTemplateString,
		IsMJML:              addWorkflowRequest.EmailTemplate.IsMJML,
		StaticAttachments: []struct {
			Content  string
			FileName string
		}{},
		TemplatedPDFs: []struct {
			TemplateName   string
			TemplateString string
		}{},
	}

	err = ah.Usecase.AddWorkflow(workflowDomain)
	if err != nil {
		handler.HandleError(w, r, http.StatusInternalServerError, fmt.Sprintf("Adding workflow with name %v failed", workflowName))
		return
	}
	resultString := fmt.Sprintf("Added workflow with name %v", workflowName)
	render.Status(r, http.StatusCreated)
	render.PlainText(w, r, resultString)
}

func (ah *APIHandler) GetWorkflowByName(w http.ResponseWriter, r *http.Request, workflowName string) {
	var workflowDomain *domain.WorkflowInfo
	var err error
	workflowDomain, err = ah.Usecase.GetWorkflowByName(workflowName)
	if err != nil {
		handler.HandleError(w, r, http.StatusInternalServerError, "Error getting workflow")
		return
	}
	if workflowDomain.Name == "" {
		handler.HandleError(w, r, http.StatusNotFound, fmt.Sprintf("Workflow with name %s not found", workflowName))
		return
	}
	dtoWorkflow := server.GetWorkflowResponse{
		Name: workflowDomain.Name,
	}

	for _, domainInput := range workflowDomain.RequiredInputs {
		dtoInput := struct {
			EmailTemplate struct {
				Placeholders *[]string `json:"placeholders,omitempty"`
				TemplateName *string   `json:"templateName,omitempty"`
			} `json:"emailTemplate"`
			PdfTemplates []struct {
				Placeholders *[]string `json:"placeholders,omitempty"`
				TemplateName *string   `json:"templateName,omitempty"`
			} `json:"pdfTemplates"`
			ToEmail string `json:"toEmail"`
			ToName  string `json:"toName"`
		}{
			ToEmail: domainInput.ToEmail,
			ToName:  domainInput.ToName,
		}

		dtoInput.EmailTemplate.TemplateName = handler.ToPointer(domainInput.EmailTemplate.TemplateName)
		dtoInput.EmailTemplate.Placeholders = handler.ToSlicePointer(domainInput.EmailTemplate.Placeholders)

		for _, domainPdfTemplate := range domainInput.PdfTemplates {
			dtoPdfTemplate := struct {
				Placeholders *[]string `json:"placeholders,omitempty"`
				TemplateName *string   `json:"templateName,omitempty"`
			}{
				TemplateName: handler.ToPointer(domainPdfTemplate.TemplateName),
				Placeholders: handler.ToSlicePointer(domainPdfTemplate.Placeholders),
			}
			// Append the mapped PDF template to the DTO input
			dtoInput.PdfTemplates = append(dtoInput.PdfTemplates, dtoPdfTemplate)
		}

		// Append the mapped RequiredInput to the DTO workflow
		dtoWorkflow.RequiredInputs = append(dtoWorkflow.RequiredInputs, dtoInput)
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, dtoWorkflow)
}

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
		EmailTemplate: struct {
			Placeholders map[string]*string
			TemplateName string
		}{
			Placeholders: handler.ConvertPlaceholders(useWorkflowRequest.EmailTemplate.Placeholders),
			TemplateName: useWorkflowRequest.EmailTemplate.TemplateName,
		},
		ToEmail: useWorkflowRequest.ToEmail,
		ToName:  useWorkflowRequest.ToName,
	}

	//TODO How to handle multiple PDF templates?
	if useWorkflowRequest.PdfTemplate != nil {
		useWorkflowRequestDomain.PdfTemplate = &struct {
			Placeholders map[string]*string
			TemplateName string
		}{
			Placeholders: handler.ConvertPlaceholders(useWorkflowRequest.PdfTemplate.Placeholders),
			TemplateName: useWorkflowRequest.PdfTemplate.TemplateName,
		}
	}

	err = ah.Usecase.UseWorkflow(useWorkflowRequestDomain)
	if err != nil {
		handler.HandleError(w, r, http.StatusInternalServerError, fmt.Sprintf("Using workflow with name %v failed", workflowName))
		return
	}
	resultString := fmt.Sprintf("Used workflow with name %v", workflowName)
	render.Status(r, http.StatusCreated)
	render.PlainText(w, r, resultString)
}
