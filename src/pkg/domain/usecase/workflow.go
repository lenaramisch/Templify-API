package usecase

import (
	"log/slog"

	domain "templify/pkg/domain/model"
)

func (u *Usecase) AddWorkflow(workflow *domain.WorkflowCreateRequest) error {
	err := u.repository.AddWorkflow(workflow)
	if err != nil {
		slog.With("workflowName", workflow.Name).Debug("Could not add workflow to repo")
		return err
	}
	return nil
}

func (u *Usecase) GetWorkflowByName(workflowName string) (*domain.Workflow, error) {
	workflow, err := u.repository.GetWorkflowByName(workflowName)
	if err != nil {
		slog.With("workflowName", workflowName).Debug("Could not get workflow from repo")
		return nil, err
	}
	return workflow, nil
}
