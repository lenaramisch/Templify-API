package usecase

import (
	"fmt"
	domain "templify/pkg/domain/model"
)

func (u *Usecase) AddEmailTemplate(r *domain.Template) error {
	err := u.repository.AddEmailTemplate(r.Name, r.TemplateStr, r.IsMJML)
	if err != nil {
		fmt.Println("=== Error ===")
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func (u *Usecase) GetEmailPlaceholders(templateName string) ([]string, error) {
	domainTemplate, err := u.repository.GetEmailTemplateByName(templateName)
	if err != nil {
		u.log.With("templateName", templateName).Debug("Could not get template from repo")
		return nil, err
	}
	return ExtractPlaceholders(domainTemplate.TemplateStr), nil
}

func (u *Usecase) GetFilledTemplateString(templateName string, values map[string]string) (string, error) {
	domainTemplate, err := u.repository.GetEmailTemplateByName(templateName)
	if err != nil {
		u.log.Debug("Error getting template by name")
		return "", err
	}

	filledTemplateString, err := FillTemplate(domainTemplate.TemplateStr, values)
	if err != nil {
		u.log.Debug("Error filling template placeholders")
		return "", err
	}

	//Check if isMJML
	if domainTemplate.IsMJML {
		// Render MJML
		filledHTMLString, err := u.mjmlService.RenderMJML(filledTemplateString)
		if err != nil {
			u.log.Debug("Error rendering mjml template")
			return "", err
		}
		return filledHTMLString, nil
	}
	return filledTemplateString, nil
}

func (u *Usecase) SendRawEmail(r *domain.EmailRequest) error {
	return u.emailSender.SendEmail(r)
}

func (u *Usecase) SendTemplatedEmail(r *domain.EmailTemplateSendRequest) error {
	// Get template by name
	templateDomain, err := u.GetEmailTemplateByName(r.TemplateName)
	if err != nil {
		u.log.Debug("Error getting template by name")
		return err
	}
	// Fill placeholders
	filledTemplate, err := FillTemplate(templateDomain.TemplateStr, r.Placeholders)
	if err != nil {
		u.log.Debug("Error filling template placeholders")
		return err
	}

	var emailRequest = domain.EmailRequest{
		ToEmail:        r.ToEmail,
		ToName:         r.ToName,
		Subject:        r.Subject,
		AttachmentInfo: r.AttachmentInfo,
	}

	// Check if isMJML
	if templateDomain.IsMJML {
		// Render MJML
		htmlString, err := u.mjmlService.RenderMJML(filledTemplate)
		if err != nil {
			u.log.Debug("Error rendering mjml template")
			return err
		}
		emailRequest.MessageBody = htmlString
	} else {
		emailRequest.MessageBody = filledTemplate
	}

	// Send email
	err = u.emailSender.SendEmail(&emailRequest)
	if err != nil {
		u.log.Debug("Error sending email")
		return err
	}
	return nil
}

func (u *Usecase) prepareMJMLBody(templateName string, values map[string]string) (*string, error) {
	// Get template and fill placeholders
	mjml, err := FillTemplate(templateName, values)
	if err != nil {
		u.log.
			With(
				"TemplateName", templateName,
				"Values", values,
				"Error", err.Error(),
			).
			Debug("Error filling mjml template with values")
		return nil, err
	}

	htmlString, err := u.mjmlService.RenderMJML(mjml)
	if err != nil {
		u.log.With(
			"MJML", mjml,
			"Error", err.Error(),
		).Debug("Error rendering mjml template")
		return nil, err
	}
	return &htmlString, nil
}

func (u *Usecase) GetEmailTemplateByName(templateName string) (*domain.Template, error) {
	templateDomain, err := u.repository.GetEmailTemplateByName(templateName)
	if err != nil {
		fmt.Println("=== Error ===")
		fmt.Println(err.Error())
		return nil, err
	}
	return templateDomain, nil
}
