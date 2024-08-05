package usecase

import (
	"fmt"
	"log/slog"
	domain "templify/pkg/domain/model"
)

func (u *Usecase) SendSMS(toNumber string, messageBody string) error {
	return u.smsSender.SendSMS(toNumber, messageBody)
}

func (u *Usecase) AddSMSTemplate(templateName string, SMSTemplString string) error {
	err := u.repository.AddSMSTemplate(templateName, SMSTemplString)
	if err != nil {
		fmt.Println("=== Error ===")
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func (u *Usecase) GetSMSTemplateByName(templateName string) (*domain.Template, error) {
	template, err := u.repository.GetSMSTemplateByName(templateName)
	if err != nil {
		slog.With("templateName", templateName).Debug("Could not get template from repo")
		return nil, err
	}
	return template, nil
}

func (u *Usecase) GetSMSPlaceholders(templateName string) ([]string, error) {
	domainTemplate, err := u.repository.GetSMSTemplateByName(templateName)
	if err != nil {
		slog.With("templateName", templateName).Debug("Could not get template from repo")
		return nil, err
	}
	return ExtractPlaceholders(domainTemplate.TemplateStr), nil
}

func (u *Usecase) GetFilledSMSTemplate(templateName string, templateFillRequest map[string]string) (string, error) {
	domainTemplate, err := u.repository.GetEmailTemplateByName(templateName)
	if err != nil {
		slog.Debug("Error getting template by name")
		return "", err
	}
	filledTemplate, err := FillTemplate(domainTemplate.TemplateStr, templateFillRequest)
	if err != nil {
		slog.Debug("Error filling template placeholders")
		return "", err
	}
	return filledTemplate, err
}
