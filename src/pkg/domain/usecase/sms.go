package usecase

import (
	"fmt"
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
		u.log.With("templateName", templateName).Debug("Could not get template from repo")
		return nil, err
	}
	return template, nil
}

func (u *Usecase) GetSMSPlaceholders(templateName string) ([]string, error) {
	domainTemplate, err := u.repository.GetSMSTemplateByName(templateName)
	if err != nil {
		u.log.With("templateName", templateName).Debug("Could not get template from repo")
		return nil, err
	}
	return ExtractPlaceholders(domainTemplate.TemplateStr), nil
}

func (u *Usecase) GetFilledSMSTemplate(templateName string, placeholders map[string]string) (string, error) {
	domainTemplate, err := u.repository.GetSMSTemplateByName(templateName)
	if err != nil {
		u.log.Debug("Error getting template by name")
		return "", err
	}
	u.log.With(
		"templateStringUnfilled", domainTemplate.TemplateStr,
	).Debug("Unfilled template string")
	filledTemplate, err := FillTemplate(domainTemplate.TemplateStr, placeholders)
	if err != nil {
		u.log.Debug("Error filling template placeholders")
		return "", err
	}
	return filledTemplate, err
}
