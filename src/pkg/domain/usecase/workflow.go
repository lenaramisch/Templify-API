package usecase

import (
	"log/slog"
	"strings"

	domain "templify/pkg/domain/model"
)

func (u *Usecase) AddWorkflow(workflow *domain.WorkflowCreateRequest) error {
	//add workflow
	err := u.repository.AddWorkflow(workflow)
	if err != nil {
		slog.With("workflowName", workflow.Name).Debug("Could not add workflow to repo")
		return err
	}
	//add email template
	err = u.repository.AddEmailTemplate(workflow.EmailTemplateName, workflow.EmailTemplateString, workflow.IsMJML)
	if err != nil {
		slog.With("workflowName", workflow.Name).Debug("Could not add email template to repo")
		return err
	}
	//add pdf templates
	for _, pdfTemplate := range workflow.TemplatedPDFs {
		err = u.repository.AddPDFTemplate(pdfTemplate.TemplateName, pdfTemplate.TemplateString)
		if err != nil {
			slog.With("workflowName", workflow.Name).Debug("Could not add pdf template to repo")
			return err
		}
	}
	//add static attachments
	for _, staticAttachment := range workflow.StaticAttachments {
		err = u.repository.SavePDF(staticAttachment.FileName, staticAttachment.Content)
		if err != nil {
			slog.With("workflowName", workflow.Name).Debug("Could not add static attachment to repo")
			return err
		}
	}
	return nil
}

func (u *Usecase) GetWorkflowByName(workflowName string) (*domain.WorkflowInfo, error) {
	workflowRaw, err := u.repository.GetWorkflowByName(workflowName)
	if err != nil {
		slog.With("workflowName", workflowName).Debug("Could not get workflow from repo")
		return nil, err
	}
	var workflowInfo = &domain.WorkflowInfo{}

	workflowInfo.RequiredInputs = append(workflowInfo.RequiredInputs, struct {
		ToEmail       string
		ToName        string
		EmailTemplate struct {
			TemplateName string
			Placeholders []string
		}
		PdfTemplates []struct {
			TemplateName string
			Placeholders []string
		}
	}{})

	// Get email template and placeholders
	emailTemplate, err := u.repository.GetEmailTemplateByName(workflowRaw.EmailTemplateName)
	if err != nil {
		slog.With("workflowName", workflowName).Debug("Could not get email template from repo")
		return nil, err
	}

	emailTemplatePlaceholders := ExtractPlaceholders(emailTemplate.TemplateStr)

	workflowInfo.RequiredInputs[0].EmailTemplate = struct {
		TemplateName string
		Placeholders []string
	}{
		TemplateName: emailTemplate.Name,
		Placeholders: emailTemplatePlaceholders,
	}

	// Split the PDF template names string into single names
	pdfTemplateNames := strings.Split(workflowRaw.PDFTemplateNames, ",")
	for _, templateName := range pdfTemplateNames {
		// Get each PDF template and placeholders
		pdfTemplate, err := u.repository.GetPDFTemplateByName(templateName)
		if err != nil {
			slog.With("workflowName", workflowName).Debug("Could not get pdf template from repo")
			return nil, err
		}

		pdfTemplatePlaceholders := ExtractPlaceholders(pdfTemplate.TemplateStr)

		// Append PDF template details to the PdfTemplates slice
		workflowInfo.RequiredInputs[0].PdfTemplates = append(workflowInfo.RequiredInputs[0].PdfTemplates, domain.PdfTemplate{
			TemplateName: templateName,
			Placeholders: pdfTemplatePlaceholders,
		})
	}

	workflowInfo.Name = workflowRaw.Name

	return workflowInfo, nil
}
