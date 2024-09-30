package usecase

import (
	domain "templify/pkg/domain/model"
)

func (u *Usecase) SendSMS(smsRequest domain.SmsRequest) error {
	return u.smsSender.SendSMS(smsRequest)
}

func (u *Usecase) AddSMSTemplate(templateName string, SMSTemplString string) error {
	err := u.repository.AddSMSTemplate(templateName, SMSTemplString)
	if err != nil {
		u.log.Debug("Error adding SMS template")
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

func (u *Usecase) GetFilledSMSTemplate(templateName string, placeholderValues map[string]string) (string, error) {
	domainTemplate, err := u.repository.GetSMSTemplateByName(templateName)
	if err != nil {
		u.log.Debug("Error getting template by name")
		return "", err
	}

	filledTemplate, err := FillTemplate(domainTemplate.TemplateStr, placeholderValues)
	if err != nil {
		u.log.Debug("Error filling template placeholders")
		return "", err
	}
	return filledTemplate, err
}

func (u *Usecase) SendTemplatedSMS(SendTemplatedSMSRequest domain.SMSTemplateFillRequest, templateName string) error {
	filledTemplate, err := u.GetFilledSMSTemplate(templateName, SendTemplatedSMSRequest.Placeholders)
	if err != nil {
		u.log.Debug("Error getting filled template")
		return err
	}

	sendSMSRequest := domain.SmsRequest{
		ToNumber:    SendTemplatedSMSRequest.ReceiverPhoneNumber,
		MessageBody: filledTemplate,
	}

	err = u.smsSender.SendSMS(sendSMSRequest)
	if err != nil {
		u.log.Debug("Error sending SMS")
		return err
	}
	return nil
}
