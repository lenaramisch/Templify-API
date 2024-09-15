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

	var workflowDomain = &domain.Workflow{
		Name:              workflowName,
		EmailTemplateName: addWorkflowRequest.EmailTemplateName,
		StaticAttachments: []struct {
			Content  string
			FileName string
		}(addWorkflowRequest.StaticAttachments),
		TemplatedPDFs: []struct {
			TemplateName   string
			TemplateString string
		}(addWorkflowRequest.TemplatedAttachments),
		EmailSubject: addWorkflowRequest.EmailSubject,
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
	var workflowDomain *domain.Workflow
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
	render.Status(r, http.StatusOK)
	render.JSON(w, r, workflowDomain)
}

func (ah *APIHandler) UseWorkflow(w http.ResponseWriter, r *http.Request, workflowName string) {

}
