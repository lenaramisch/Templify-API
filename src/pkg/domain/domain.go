package domain

import (
	_ "embed"
)

//go:embed template_test.mjml
var embeddedTemplate string

type EmailSender interface {
	SendEmail(toEmail string, toName string, subject string, message string) error
}

type SMSSender interface {
	SendSMS(toNumber string, messageBody string) error
}

type MJMLService interface {
	GetTemplatePlaceholders(templateName string) ([]string, error)
	TemplatePostRequest(templateName string) error
	FillTemplatePlaceholders(templateName string, values map[string]interface{}) (string, error)
}

type Usecase struct {
	emailSender EmailSender
	smsSender   SMSSender
	mjmlService MJMLService
}

func NewUsecase(emailsender EmailSender, smsSender SMSSender, mjmlService MJMLService) *Usecase {
	return &Usecase{
		emailSender: emailsender,
		smsSender:   smsSender,
		mjmlService: mjmlService,
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
	// return u.mjmlService.GetTemplatePlaceholders(templateName) // TODO THIS IS CORRECT
	return u.mjmlService.GetTemplatePlaceholders(embeddedTemplate)
}

func (u *Usecase) FillTemplatePlaceholders(templateName string, values map[string]any) (string, error) {
	return u.mjmlService.FillTemplatePlaceholders(templateName, values)
}
