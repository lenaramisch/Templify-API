package usecase

import (
	"fmt"
	"log/slog"
	domain "templify/pkg/domain/model"
)

func (u *Usecase) AddEmailTemplate(templateName string, MJMLString string) error {
	err := u.repository.AddEmailTemplate(templateName, MJMLString)
	if err != nil {
		fmt.Println("=== Error ===")
		fmt.Println(err.Error())
		return err
	}
	return nil
}

// TODO rename all funcs that have template to email template
// also reuse util func and add logging?
func (u *Usecase) GetEmailPlaceholders(templateName string) ([]string, error) {
	domainTemplate, err := u.repository.GetEmailTemplateByName(templateName)
	if err != nil {
		slog.With("templateName", templateName).Debug("Could not get template from repo")
		return nil, err
	}
	return ExtractPlaceholders(domainTemplate.TemplateStr), nil
}

func (u *Usecase) GetFilledMJMLTemplate(templateName string, values map[string]string) (string, error) {
	domainTemplate, err := u.repository.GetEmailTemplateByName(templateName)
	if err != nil {
		slog.Debug("Error getting template by name")
		return "", err
	}
	filledTemplate, err := FillTemplate(domainTemplate.TemplateStr, values)
	if err != nil {
		slog.Debug("Error filling template placeholders")
		return "", err
	}
	return filledTemplate, err
}

func (u *Usecase) SendRawEmail(r *domain.EmailRequest) error {
	return u.emailSender.SendEmail(r)
}

func (u *Usecase) SendMJMLEmail(r *domain.EmailRequest, templateName string, values map[string]string) error {
	emailBody, err := u.prepareMJMLBody(templateName, values)
	if err != nil {
		return err
	}
	r.MessageBody = *emailBody
	err = u.emailSender.SendEmail(r)
	if err != nil {
		slog.Debug("Error sending email without attachment")
		return err
	}
	return nil
}

func (u *Usecase) prepareMJMLBody(templateName string, values map[string]string) (*string, error) {
	// Get template and fill placeholders
	mjml, err := u.GetFilledMJMLTemplate(templateName, values)
	if err != nil {
		slog.
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
		slog.With(
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
