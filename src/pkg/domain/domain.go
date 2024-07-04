package domain

import (
	_ "embed"
	"fmt"
	"log/slog"
)

type EmailSender interface {
	SendEmail(toEmail string, toName string, subject string, message string) error
}

type SMSSender interface {
	SendSMS(toNumber string, messageBody string) error
}

type MJMLService interface {
	GetTemplatePlaceholders(template Template) ([]string, error)
	FillTemplatePlaceholders(domainTemplate Template, values map[string]any) (string, error)
	RenderMJML(MJMLString string) (string, error)
}

type Repository interface {
	GetTemplateByName(name string) (*Template, error)
	AddTemplate(name string, mjmlString string) error
}

type Usecase struct {
	emailSender EmailSender
	smsSender   SMSSender
	mjmlService MJMLService
	repository  Repository
}

func NewUsecase(emailsender EmailSender, smsSender SMSSender, mjmlService MJMLService, repository Repository) *Usecase {
	return &Usecase{
		emailSender: emailsender,
		smsSender:   smsSender,
		mjmlService: mjmlService,
		repository:  repository,
	}
}

// domain layer functions (usecases with actual business logic)
// Here we currently don't have any logic and forward to our services
// Later we may have sequential service calls or mapping, some logic etc.

func (u *Usecase) SendEmail(toEmail string, toName string, subject string, message string) error {
	return u.emailSender.SendEmail(toEmail, toName, subject, message)
}

func (u *Usecase) SendSMS(toNumber string, messageBody string) error {
	return u.smsSender.SendSMS(toNumber, messageBody)
}

func (u *Usecase) GetTemplatePlaceholders(templateName string) ([]string, error) {
	domainTemplate, err := u.repository.GetTemplateByName(templateName)
	if err != nil {
		return nil, err
	}
	return u.mjmlService.GetTemplatePlaceholders(*domainTemplate)
}

func (u *Usecase) FillTemplatePlaceholders(templateName string, shouldBeSent bool, toEmail string, toName string, subject string, values map[string]any) (string, error) {
	domainTemplate, err := u.repository.GetTemplateByName(templateName)
	if err != nil {
		slog.Debug("Error getting template by name")
		return "", err
	}
	filledTemplate, err := u.mjmlService.FillTemplatePlaceholders(*domainTemplate, values)
	if err != nil {
		slog.Debug("Error filling template placeholders")
		return "", err
	}
	if !shouldBeSent {
		return filledTemplate, nil
	}
	//TODO send MJMLString to MJML Service to get html string
	htmlString, err := u.mjmlService.RenderMJML(filledTemplate)
	if err != nil {
		slog.Debug("Error rendering mjml template")
		return "", err
	}
	//TODO send html string as email with sendgrid
	err = u.emailSender.SendEmail(toEmail, toName, subject, htmlString)
	if err != nil {
		slog.Debug("Error sending email")
		return "", err
	}
	return htmlString, nil
}

func (u *Usecase) AddTemplate(templateName string, MJMLString string) error {
	err := u.repository.AddTemplate(templateName, MJMLString)
	if err != nil {
		fmt.Println("=== Error ===")
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func (u *Usecase) GetTemplateByName(templateName string) (*Template, error) {
	templateDomain, err := u.repository.GetTemplateByName(templateName)
	if err != nil {
		fmt.Println("=== Error ===")
		fmt.Println(err.Error())
		return nil, err
	}
	return templateDomain, nil
}
