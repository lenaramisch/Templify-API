package usecase

import (
	"context"
	"strings"
	domain "templify/pkg/domain/model"
)

func (u *Usecase) AddWorkflow(ctx context.Context, workflow *domain.WorkflowCreateRequest) error {
	err := u.repository.AddWorkflow(ctx, workflow)
	if err != nil {
		return err
	}
	return nil
}

func (u *Usecase) GetWorkflowByName(ctx context.Context, workflowName string) (*domain.GetWorkflowResponse, error) {
	workflowRaw, err := u.repository.GetWorkflowByName(ctx, workflowName)
	if err != nil {
		return nil, err
	}
	var getWorkflowResponse = &domain.GetWorkflowResponse{}

	getWorkflowResponse.Name = workflowRaw.Name
	getWorkflowResponse.EmailSubject = workflowRaw.EmailSubject

	// Get email template and placeholders
	emailTemplate, err := u.repository.GetEmailTemplateByName(ctx, workflowRaw.EmailTemplateName)
	if err != nil {
		u.log.With("workflowName", workflowName).Debug("Could not get email template from repo")
		return nil, err
	}

	emailTemplatePlaceholders := ExtractPlaceholders(emailTemplate.TemplateStr)

	getWorkflowResponse.EmailTemplate = domain.TemplateInfo{
		TemplateName: emailTemplate.Name,
		Placeholders: emailTemplatePlaceholders,
	}

	// Split the PDF template names string into single names
	var pdfTemplateNames []string
	pdfTemplateNames = append(pdfTemplateNames, strings.Split(workflowRaw.PDFTemplateNames, ",")...)

	for _, templateName := range pdfTemplateNames {
		// Get each PDF template and placeholders
		pdfTemplate, err := u.repository.GetPDFTemplateByName(ctx, templateName)
		if err != nil {
			u.log.With("workflowName", workflowName).Debug("Could not get pdf template from repo")
			return nil, err
		}

		pdfTemplatePlaceholders := ExtractPlaceholders(pdfTemplate.TemplateStr)

		getWorkflowResponse.PDFTemplates = append(getWorkflowResponse.PDFTemplates, domain.TemplateInfo{
			TemplateName: templateName,
			Placeholders: pdfTemplatePlaceholders,
		})
	}

	// Split the static attachment names string into single names
	getWorkflowResponse.StaticAttachments = append(getWorkflowResponse.StaticAttachments, strings.Split(workflowRaw.StaticAttachments, ",")...)

	return getWorkflowResponse, nil
}

func (u *Usecase) UseWorkflow(ctx context.Context, workflowUseRequest *domain.WorkflowUseRequest) error {
	// get pdf templates by name
	var pdfTemplates []domain.Template
	var pdfTemplatePlaceholders []map[string]string
	for _, pdfTemplate := range workflowUseRequest.PdfTemplates {
		pdfTemplatePlaceholders = append(pdfTemplatePlaceholders, pdfTemplate.Placeholders)
		pdfTemplate, err := u.repository.GetPDFTemplateByName(ctx, pdfTemplate.TemplateName)
		if err != nil {
			return err
		}
		pdfTemplates = append(pdfTemplates, *pdfTemplate)
	}

	// fill pdf template
	// range over the pdf templates and fill them
	var filledPDFTemplates []string
	for i, pdfTemplate := range pdfTemplates {
		filledPdfTemplate, err := FillTemplate(pdfTemplate.TemplateStr, pdfTemplatePlaceholders[i])
		if err != nil {
			return err
		}
		filledPDFTemplates = append(filledPDFTemplates, filledPdfTemplate)
	}

	attachmentData := []domain.AttachmentInfo{}

	for i, pdfTemplate := range pdfTemplates {
		filledTemplate := filledPDFTemplates[i]
		filledPdfFile, err := u.typstService.RenderTypst(filledTemplate)
		if err != nil {
			return err
		}
		attachmentData = append(attachmentData, domain.AttachmentInfo{
			FileName:      pdfTemplate.Name,
			FileExtension: "pdf",
			FileBytes:     filledPdfFile,
		})
	}

	// get workflow
	workflowInfo, err := u.GetWorkflowByName(ctx, workflowUseRequest.Name)
	if err != nil {
		return err
	}

	var staticAttachments []domain.StaticFile
	for _, attachmentName := range workflowInfo.StaticAttachments {
		// split file extrension from name
		splitString := strings.SplitN(attachmentName, ".", 2)
		var fileName, extension string
		if len(splitString) == 2 {
			fileName = splitString[0]
			extension = splitString[1]
		} else {
			return domain.ErrorAttachmentNameInvalid{AttachmentName: attachmentName}
		}
		downloadFileRequest := domain.FileDownloadRequest{
			FileName:  fileName,
			Extension: extension,
		}
		file, err := u.filemanagerService.DownloadFile(downloadFileRequest)
		if err != nil {
			return err
		}

		//append attachment to staticAttachments
		staticAttachments = append(staticAttachments, domain.StaticFile{
			FileName:  splitString[0],
			Extension: splitString[1],
			Content:   file,
		})
	}

	for _, attachment := range staticAttachments {
		attachmentData = append(attachmentData, domain.AttachmentInfo{
			FileName:      attachment.FileName,
			FileExtension: attachment.Extension,
			FileBytes:     attachment.Content,
		})
	}

	emailRequest := &domain.EmailTemplateSendRequest{
		ToEmail:        workflowUseRequest.ToEmail,
		ToName:         workflowUseRequest.ToName,
		Subject:        workflowInfo.EmailSubject,
		TemplateName:   workflowUseRequest.EmailTemplate.TemplateName,
		Placeholders:   workflowUseRequest.EmailTemplate.Placeholders,
		AttachmentInfo: attachmentData,
	}

	u.log.With("emailRequest", emailRequest).Debug("Email Request")

	// send email
	err = u.SendTemplatedEmail(ctx, emailRequest)
	if err != nil {
		return err
	}
	return nil
}
