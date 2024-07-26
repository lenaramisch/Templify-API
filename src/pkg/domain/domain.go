package domain

import (
	_ "embed"
)

type EmailSender interface {
	SendEmail(domainEmailReq *EmailRequest) error
}

type SMSSender interface {
	SendSMS(toNumber string, messageBody string) error
}

type MJMLService interface {
	RenderMJML(MJMLString string) (string, error)
}

type Repository interface {
	// Email
	GetEmailTemplateByName(name string) (*Template, error)
	AddEmailTemplate(name string, mjmlString string) error
	// PDF
	GetPDFTemplateByName(name string) (*Template, error)
	AddPDFTemplate(name string, typstString string) error
}

type TypstService interface {
	RenderTypst(typstString string) ([]byte, error)
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
