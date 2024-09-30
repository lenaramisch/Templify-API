package usecase

import (
	"strings"
	domain "templify/pkg/domain/model"
)

func (u *Usecase) AddWorkflow(workflow *domain.WorkflowCreateRequest) error {
	//add workflow
	err := u.repository.AddWorkflow(workflow)
	if err != nil {
		u.log.With("workflowName", workflow.Name).Debug("Could not add workflow to repo")
		return err
	}
	return nil
}

func (u *Usecase) GetWorkflowByName(workflowName string) (*domain.GetWorkflowResponse, error) {
	workflowRaw, err := u.repository.GetWorkflowByName(workflowName)
	if err != nil {
		u.log.With("workflowName", workflowName).Debug("Could not get workflow from repo")
		return nil, err
	}
	var getWorkflowResponse = &domain.GetWorkflowResponse{}

	getWorkflowResponse.Name = workflowRaw.Name
	getWorkflowResponse.EmailSubject = workflowRaw.EmailSubject

	// Get email template and placeholders
	emailTemplate, err := u.repository.GetEmailTemplateByName(workflowRaw.EmailTemplateName)
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
	for _, name := range strings.Split(workflowRaw.PDFTemplateNames, ",") {
		pdfTemplateNames = append(pdfTemplateNames, name)
	}

	for _, templateName := range pdfTemplateNames {
		// Get each PDF template and placeholders
		pdfTemplate, err := u.repository.GetPDFTemplateByName(templateName)
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
	for _, name := range strings.Split(workflowRaw.StaticAttachments, ",") {
		getWorkflowResponse.StaticAttachments = append(getWorkflowResponse.StaticAttachments, name)
	}

	return getWorkflowResponse, nil
}

func (u *Usecase) UseWorkflow(workflowUseRequest *domain.WorkflowUseRequest) error {
	// get pdf template by name
	pdfTemplate, err := u.repository.GetPDFTemplateByName(workflowUseRequest.PdfTemplate.TemplateName)
	if err != nil {
		u.log.With("templateName", workflowUseRequest.PdfTemplate.TemplateName).Debug("Could not get pdf template from repo")
		return err
	}

	// fill pdf template
	filledPdfTemplate, err := FillTemplate(pdfTemplate.TemplateStr, workflowUseRequest.PdfTemplate.Placeholders)
	if err != nil {
		u.log.With("templateName", workflowUseRequest.PdfTemplate.TemplateName).Debug("Could not fill pdf template")
		return err
	}

	splitString := strings.SplitN(workflowUseRequest.PdfTemplate.TemplateName, ".", 2)

	var fileName, extension string
	if len(splitString) == 2 {
		fileName = splitString[0]
		extension = splitString[1]
	} else {
		fileName = workflowUseRequest.PdfTemplate.TemplateName
		extension = "pdf"
	}

	attachmentData := []domain.AttachmentInfo{}

	pdfAttachment := domain.StaticFile{
		FileName:  fileName,
		Extension: extension,
	}

	filledPdfFile, err := u.typstService.RenderTypst(filledPdfTemplate)
	attachmentData = append(attachmentData, domain.AttachmentInfo{
		FileName:      pdfAttachment.FileName,
		FileExtension: pdfAttachment.Extension,
		FileBytes:     filledPdfFile,
	})

	// get workflow
	workflowInfo, err := u.GetWorkflowByName(workflowUseRequest.Name)
	if err != nil {
		u.log.With("workflowName", workflowUseRequest.Name).Debug("Could not get workflow from repo")
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
			u.log.With("attachmentName", attachmentName).Debug("Static file name does not contain an extension")
			return err
		}
		downloadFileRequest := domain.FileDownloadRequest{
			FileName:  fileName,
			Extension: extension,
		}
		file, err := u.fileManagerService.DownloadFile(downloadFileRequest)
		if err != nil {
			u.log.With("attachmentName", attachmentName).Debug("Downloading file failed")
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
	err = u.SendTemplatedEmail(emailRequest)
	if err != nil {
		u.log.With("emailRequest", emailRequest).Debug("Could not send email")
		return err
	}
	return nil
}
