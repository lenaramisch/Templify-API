package domain

import (
	_ "embed"
	"fmt"
	"log/slog"
)

type EmailSender interface {
	SendEmail(emailRequest *EmailRequest) error
	SendEmailWithAttachment(message string, domainEmailReq *EmailRequestAttm) error
}

type SMSSender interface {
	SendSMS(toNumber string, messageBody string) error
}

type MJMLService interface {
	GetTemplatePlaceholders(template Template) ([]string, error)
	FillTemplatePlaceholders(domainTemplate Template, values map[string]string) (string, error)
	RenderMJML(MJMLString string) (string, error)
}

type Repository interface {
	GetTemplateByName(name string) (*Template, error)
	AddEmailTemplate(name string, mjmlString string) error
	AddPDFTemplate(name string, typstString string) error
	GetPDFTemplateByName(name string) (*PDFTemplate, error)
}

type TypstService interface {
	FillPDFTemplatePlaceholders(typstTempl *PDFTemplate, placeholders map[string]string) (string, error)
	GetPDFTemplatePlaceholders(typstString string) ([]string, error)
}

type Usecase struct {
	emailSender  EmailSender
	smsSender    SMSSender
	mjmlService  MJMLService
	repository   Repository
	typstService TypstService
}

func NewUsecase(emailSender EmailSender, smsSender SMSSender, mjmlService MJMLService, repository Repository, typstService TypstService) *Usecase {
	return &Usecase{
		emailSender:  emailSender,
		smsSender:    smsSender,
		mjmlService:  mjmlService,
		repository:   repository,
		typstService: typstService,
	}
}

func (u *Usecase) SendEmail(emailRequest *EmailRequest) error {
	return u.emailSender.SendEmail(emailRequest)
}

func (u *Usecase) SendEmailWithAttachment(message string, domainEmailReq *EmailRequestAttm) error {
	return u.emailSender.SendEmailWithAttachment(message, domainEmailReq)
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

func (u *Usecase) FillTemplatePlaceholders(templateName string, domainFillTempl *TemplateFillRequest) (string, error) {
	domainTemplate, err := u.repository.GetTemplateByName(templateName)
	if err != nil {
		slog.Debug("Error getting template by name")
		return "", err
	}
	filledTemplate, err := u.mjmlService.FillTemplatePlaceholders(*domainTemplate, domainFillTempl.Placeholders)
	if err != nil {
		slog.Debug("Error filling template placeholders")
		return "", err
	}
	if !domainFillTempl.ShouldBeSent {
		return filledTemplate, nil
	}

	htmlString, err := u.mjmlService.RenderMJML(filledTemplate)
	if err != nil {
		slog.Debug("Error rendering mjml template")
		return "", err
	}
	var emailRequest EmailRequest
	emailRequest.MessageBody = htmlString
	emailRequest.Subject = domainFillTempl.Subject
	emailRequest.ToEmail = domainFillTempl.ToEmail
	emailRequest.ToName = domainFillTempl.ToName
	err = u.emailSender.SendEmail(&emailRequest)
	if err != nil {
		slog.Debug("Error sending email")
		return "", err
	}
	return htmlString, nil
}

func (u *Usecase) AddEmailTemplate(templateName string, MJMLString string) error {
	err := u.repository.AddEmailTemplate(templateName, MJMLString)
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

func (u *Usecase) FillTemplatePlaceholdersAttm(templateName string, domainEmailReq *EmailRequestAttm) (string, error) {
	// Get domainTempl
	domainTemplate, err := u.repository.GetTemplateByName(templateName)
	if err != nil {
		slog.Debug("Error getting template by name")
		return "", err
	}
	// Fill templ
	filledTemplate, err := u.mjmlService.FillTemplatePlaceholders(*domainTemplate, domainEmailReq.Placeholders)
	if err != nil {
		slog.Debug("Error filling template placeholders")
		return "", err
	}
	if !domainEmailReq.ShouldBeSent {
		return filledTemplate, nil
	}

	// Render MJML string to get html
	htmlString, err := u.mjmlService.RenderMJML(filledTemplate)
	if err != nil {
		slog.Debug("Error rendering mjml template")
		return "", err
	}

	err = u.emailSender.SendEmailWithAttachment(htmlString, domainEmailReq)
	if err != nil {
		slog.Debug("Error sending email")
		return "", err
	}
	return htmlString, nil
}

func (u *Usecase) AddPDFTemplate(templateName string, typstString string) error {
	err := u.repository.AddPDFTemplate(templateName, typstString)
	if err != nil {
		fmt.Println("=== Error ===")
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func (u *Usecase) GetPDFTemplateByName(templateName string) (*PDFTemplate, error) {
	templateDomain, err := u.repository.GetPDFTemplateByName(templateName)
	if err != nil {
		fmt.Println("=== Error ===")
		fmt.Println(err.Error())
		return nil, err
	}
	return templateDomain, nil
}

func (u *Usecase) FillPDFTemplatePlaceholders(templateName string, pdfFillRequest *PDFTemplateFillRequest) (string, error) {
	pdfTempl, err := u.GetPDFTemplateByName(templateName)
	if err != nil {
		return "Getting template from db failed", err
	}
	placeholders := pdfFillRequest.Placeholders
	u.typstService.FillPDFTemplatePlaceholders(pdfTempl, placeholders)
	return "", nil
}
