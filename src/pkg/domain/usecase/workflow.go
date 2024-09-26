package usecase

import (
	b64 "encoding/base64"
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
	//add email template
	err = u.repository.AddEmailTemplate(workflow.EmailTemplateName, workflow.EmailTemplateString, workflow.IsMJML)
	if err != nil {
		u.log.With("workflowName", workflow.Name).Debug("Could not add email template to repo")
		return err
	}
	//add pdf templates
	for _, pdfTemplate := range workflow.TemplatedPDFs {
		err = u.repository.AddPDFTemplate(pdfTemplate.Name, pdfTemplate.TemplateStr)
		if err != nil {
			u.log.With("workflowName", workflow.Name).Debug("Could not add pdf template to repo")
			return err
		}
	}
	//add static attachments
	//TODO exchange SavePDF func with UploadFile to S3 Bucket (Minio)
	for _, staticAttachment := range workflow.StaticAttachments {
		err = u.repository.SavePDF(staticAttachment.FileName, staticAttachment.Content)
		if err != nil {
			u.log.With("workflowName", workflow.Name).Debug("Could not add static attachment to repo")
			return err
		}
	}
	return nil
}

func (u *Usecase) GetWorkflowByName(workflowName string) (*domain.WorkflowInfo, error) {
	workflowRaw, err := u.repository.GetWorkflowByName(workflowName)
	if err != nil {
		u.log.With("workflowName", workflowName).Debug("Could not get workflow from repo")
		return nil, err
	}
	var workflowInfo = &domain.WorkflowInfo{}

	workflowInfo.RequiredInputs = append(workflowInfo.RequiredInputs, struct {
		ToEmail           string
		ToName            string
		EmailTemplate     domain.TemplateInfo
		PdfTemplates      []domain.TemplateInfo
		StaticAttachments []string
	}{})

	// Get email template and placeholders
	emailTemplate, err := u.repository.GetEmailTemplateByName(workflowRaw.EmailTemplateName)
	if err != nil {
		u.log.With("workflowName", workflowName).Debug("Could not get email template from repo")
		return nil, err
	}

	emailTemplatePlaceholders := ExtractPlaceholders(emailTemplate.TemplateStr)

	workflowInfo.RequiredInputs[0].EmailTemplate = domain.TemplateInfo{
		TemplateName: workflowRaw.EmailTemplateName,
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

		workflowInfo.RequiredInputs[0].PdfTemplates = append(workflowInfo.RequiredInputs[0].PdfTemplates, domain.TemplateInfo{
			TemplateName: templateName,
			Placeholders: pdfTemplatePlaceholders,
		})
	}

	// Split the static attachment names string into single names
	workflowInfo.RequiredInputs[0].StaticAttachments = []string{}
	for _, name := range strings.Split(workflowRaw.StaticAttachments, ",") {
		workflowInfo.RequiredInputs[0].StaticAttachments = append(workflowInfo.RequiredInputs[0].StaticAttachments, name)
	}

	workflowInfo.Name = workflowRaw.Name
	workflowInfo.EmailSubject = workflowRaw.EmailSubject

	return workflowInfo, nil
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

	pdfAttachment := domain.PDF{
		FileName:      fileName,
		FileExtension: extension,
		Content:       filledPdfTemplate,
	}

	base64Content := b64.StdEncoding.EncodeToString([]byte(pdfAttachment.Content))
	attachmentData = append(attachmentData, domain.AttachmentInfo{
		FileName:      pdfAttachment.FileName,
		FileExtension: pdfAttachment.FileExtension,
		Base64Content: base64Content,
	})

	// get workflow
	workflowInfo, err := u.GetWorkflowByName(workflowUseRequest.Name)
	if err != nil {
		u.log.With("workflowName", workflowUseRequest.Name).Debug("Could not get workflow from repo")
		return err
	}

	//TODO exchange GetPDF func with DownloadFile from S3 Bucket (Minio)
	var staticAttachments []domain.PDF
	for _, attachmentName := range workflowInfo.RequiredInputs[0].StaticAttachments {
		content, err := u.repository.GetPDF(attachmentName)
		if err != nil {
			u.log.With("attachmentName", attachmentName).Debug("Could not get PDF from repo")
			return err
		}

		splitString = strings.SplitN(attachmentName, ".", 2)

		//append attachment to staticAttachments
		staticAttachments = append(staticAttachments, domain.PDF{
			FileName:      splitString[0],
			FileExtension: splitString[1],
			Content:       content,
		})
	}

	for _, attachment := range staticAttachments {
		attachmentData = append(attachmentData, domain.AttachmentInfo{
			FileName:      attachment.FileName,
			FileExtension: attachment.FileExtension,
			Base64Content: attachment.Content,
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
