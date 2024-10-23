package usecase

import (
	"context"
	domain "templify/pkg/domain/model"
)

func (u *Usecase) SendSMS(smsRequest domain.SmsRequest) error {
	return u.smsSender.SendSMS(smsRequest)
}

func (u *Usecase) AddSMSTemplate(ctx context.Context, template *domain.Template) error {
	err := u.repository.AddSMSTemplate(ctx, template)
	if err != nil {
		u.log.Debug("Error adding SMS template")
		return err
	}
	return nil
}

func (u *Usecase) GetSMSTemplateByName(ctx context.Context, templateName string) (*domain.Template, error) {
	template, err := u.repository.GetSMSTemplateByName(ctx, templateName)
	if err != nil {
		u.log.With("templateName", templateName).Debug("Could not get template from repo")
		return nil, err
	}
	return template, nil
}

func (u *Usecase) GetSMSPlaceholders(ctx context.Context, templateName string) ([]string, error) {
	domainTemplate, err := u.repository.GetSMSTemplateByName(ctx, templateName)
	if err != nil {
		u.log.With("templateName", templateName).Debug("Could not get template from repo")
		return nil, err
	}
	return ExtractPlaceholders(domainTemplate.TemplateStr), nil
}

func (u *Usecase) GetFilledSMSTemplate(ctx context.Context, templateName string, placeholderValues map[string]string) (string, error) {
	domainTemplate, err := u.repository.GetSMSTemplateByName(ctx, templateName)
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

func (u *Usecase) SendTemplatedSMS(ctx context.Context, SendTemplatedSMSRequest domain.SMSTemplateFillRequest, templateName string) error {
	filledTemplate, err := u.GetFilledSMSTemplate(ctx, templateName, SendTemplatedSMSRequest.Placeholders)
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
